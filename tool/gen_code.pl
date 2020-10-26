#!/usr/bin/env perl

use File::Basename;
use Data::Dumper;

my $base_dir = dirname(__FILE__) . '/..';

my $struct_info = struct_info();

process_file($base_dir . '/eventlog/fill.go');
process_file($base_dir . '/service/handlers/validate.go');

sub process_file {
    my ($file) = @_;

    open my $fh, '<:utf8', $file or die "open source file $file fail\n";
    my $content = do {
        local $/ = <$fh>;
    };
    close $fh;

    my @tpls = ();

    $content =~ s{
        (^[\ \t]*//\s*TPL\.(\w+)\.START(.*?)$)
        (.+?)
        (^[\ \t]*//\s*TPL\.\g2.END.*?$)
    }{
        my $tpl = main->can('tpl_' . lc($2));

        unless ($tpl) {
            die "tpl $2 not defined";
        }

        my $tpl_start = $1;
        my $params = $3;
        my $tpl_end = $5;

        push @tpls, $2;

        $params =~ s/(^\s*|\s*$)//g;

        my $replace = "$tpl_start\n". $tpl->($params) . "\n$tpl_end";

        $replace;
    }xmsieg;

    if (@tpls == 0) {
        print "\033[0;31m$file do not contains TPL defination\033[0m\n";
        return;
    }

    open my $fh, '>:utf8', $file or die "open source file $file for writing fail\n";
    print $fh $content;
    close $fh;

    print "Done $file : \033[0;32m", join(",", @tpls), "\033[0m\n";
}

sub tpl_check_required {
    my ($params) = @_;

    my $fields_def = $struct_info->{EventLog};
    my @fields = split /\s*,\s*/, $params;

    my @content = ();
    for my $field (@fields) {
        $field =~ s/(^\s+|\s+$)//g;
        next unless $field && $fields_def->{$field};

        my $type = $fields_def->{$field}{type};
        my $proto_type = $fields_def->{$field}{proto_type};

        if ($type eq 'string') {
            push @content, <<EOT;
	if log.$field == "" {
		return errors.Wrap(errors.ErrFieldRequired, "$field")
	}
EOT
        }

        if ($type =~ /^int(\d+)?$/ || $proto_type eq 'varint') {
            push @content, <<EOT;
	if log.$field == 0 {
		return errors.Wrap(errors.ErrFieldRequired, "$field")
	}
EOT
        }
    }
    return join("\n", @content)
}

sub tpl_event_log_fill {
    my $common_fields_def = $struct_info->{EventLogCommon};

    my @content = ();
    for my $field (sort keys %{$common_fields_def}) {
        my $type = $common_fields_def->{$field}{type};
        my $proto_type = $common_fields_def->{$field}{proto_type};

        if ($type eq 'string') {
            push @content, <<EOT;
		if log.$field == "" && common != nil && common.$field != "" {
			log.$field = common.$field
		}
EOT
        }

        if ($type =~ /^int(\d+)?$/ || $proto_type eq 'varint') {
            push @content, <<EOT;
		if log.$field == 0 && common != nil && common.$field != 0 {
			log.$field = common.$field
		}
EOT
        }
    }

    return join("\n", @content);
}

sub struct_info {
    my $file = $base_dir . '/interface-defs/event_log.pb.go';
    open my $fh, '<:utf8', $file or die "open file $file fail\n";
    my $content = do {
        local $/ = <$fh>;
    };
    close $fh;

    # event log
    my $event_log_def = get_struct_def(EventLog => $content);

    unless ($event_log_def) {
        die "Parse event log def fail\n";
    }

    # LogId      string `protobuf:"bytes,1,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
    my $event_log_fields_def = get_fields_def($event_log_def);

    # event log common
    my $event_log_common_def = get_struct_def(EventLogCommon => $content);

    unless ($event_log_common_def) {
        die "Parse event log common fail\n";
    }

    my $event_log_common_fields_def = get_fields_def($event_log_common_def);

    return {
        'EventLog' => $event_log_fields_def,
        'EventLogCommon' => $event_log_common_fields_def,
    };
}

sub get_struct_def {
    my ($struct_name, $content) = @_;

    my ($struct_def) = $content =~ m/
        type\s+${struct_name}\s+struct\s+\{
        (.+?)
        \}
    /xsm;

    return $struct_def;
}

sub get_fields_def {
    my ($struct_def) = @_;

    my %fields_def = ();

    while($struct_def =~ /^\s+(\w+)\s+(\S+)\s+`protobuf:"(\w+)/mg) {
        $fields_def{$1} = { type => $2, proto_type => $3 };
    }

    return \%fields_def;
}

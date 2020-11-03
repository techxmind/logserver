#!/usr/bin/env perl
#===============================================================================
#
#         FILE: csv_helper.pl
#
#        USAGE: ./csv_helper.pl
#
#  DESCRIPTION:
#
#      OPTIONS: ---
# REQUIREMENTS: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: YOUR NAME (),
# ORGANIZATION:
#      VERSION: 1.0
#      CREATED: 2020/11/03 20时00分21秒
#     REVISION: ---
#===============================================================================

use strict;
use warnings;
use utf8;

open my $fh, './event_log.sql' or die "$@";

my $sql = do {
    local $/ = <$fh>;
};

my @fields = $sql =~ m/^\s*(\w+)\s*(?:int|bigint|varchar|time)/img;

# shift id
shift @fields;

my $fields = join(",", @fields);

# replace logged_time => logged_time_str to get string inteadof unix timestamp
(my $csv_headers = $fields) =~ s/logged_time/logged_time_str/;

print <<EOT;
CSV Headers:
------------------------------------------------
$csv_headers

Load Data SQL:
------------------------------------------------
LOAD DATA LOCAL INFILE '{FILE}' REPLACE INTO TABLE {TABLE}
  CHARACTER SET utf8mb4
  FIELDS
    TERMINATED BY ','
    OPTIONALLY ENCLOSED BY '"'
    ESCAPED BY '"'
  LINES
    TERMINATED BY '\\n'
(
  $fields
);
EOT

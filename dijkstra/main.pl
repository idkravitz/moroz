#!/usr/bin/env perl
use strict;
use warnings;

sub usage {
    print "Usage: dijkstra.pl [option]\n" .
        "Reads expression from stdin\n" .
        " and writes to stdout\n" .
        " -h, --help:\t\t\tPrint this message\n";
}

if (@ARGV == 1) {
    if ($ARGV[0] =~ /-h|--help/) {
        usage
    } else {
        print "Error: Unrecognized option\n";
    }
    exit 0;
} elsif (@ARGV > 1) {
    print "Error: Too many arguments\n";
    exit 0;
}

my %tokenizer = (
    buffer => ""
);

sub getChar {
    $tokenizer{buffer} . <STDIN> if @tokenizer{buffer} == 0;
    return shift $

}

sub getNextToken {
    print "buf: " . $tokenizer{buffer};
    $tokenizer{buffer} = "123";
    ch = getChar 
}

getNextToken;
print "buf: " . $tokenizer{buffer};

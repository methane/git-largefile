#!/usr/bin/env python
# -*- coding: utf-8 -*-
from __future__ import print_function

import os, sys, path, boto, hashlib
import ConfigParser

BASE_DIR = path.path('~/.gitasset').expanduser()
DATA_DIR = BASE_DIR / 'data'

def get_config():
    parser = ConfigParser.SafeConfigParser()
    parser.read([BASE_DIR / 'gits3.ini'])
    return parser

def get_cache_path(hexdigest):
    return DATA_DIR / hexdigest[:2] / hexdigest[2:4] / hexdigest[4:]

def read_stdin():
    return os.fdopen(sys.stdin.fileno(), 'rb').read()

def write_stdout(content):
    return os.fdopen(sys.stdout.fileno(), 'wb').write(content)

def get_key(hexdigest):
    conf = get_config()
    s3key = conf.get('DEFAULT', 'awskey').split(':')
    bucket_name = conf.get('DEFAULT', 'bucket')
    s3 = boto.connect_s3(*s3key)
    bucket = s3.get_bucket(bucket_name)
    return boto.s3.key.Key(bucket, hexdigest)

def store():
    content = read_stdin()
    hexdigest = hashlib.sha1(content).hexdigest()
    cache_path = get_cache_path(hexdigest)
    if not cache_path.exists():
        cache_path.dirname().makedirs_p()
        cache_path.write_bytes(content)

        key = get_key(hexdigest)
        if not key.exists():
            key.set_contents_from_string(content)

    write_stdout(hexdigest)

def load():
    hexdigest = read_stdin()
    cache_path = get_cache_path(hexdigest)
    if cache_path.exists():
        contents = cache_path.bytes()
    else:
        contents = get_key(hexdigest).get_contents_as_string()
    write_stdout(contents)

if __name__ == '__main__':
    if sys.argv[1] == 'store':
        store()
    elif sys.argv[1] == 'load':
        load()
    else:
        print("usate: gits3.py <store|load>", file=sys.stderr)
        sys.exit(1)

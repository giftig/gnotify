#!/usr/bin/env python2.7

import argparse
import os
import sys

import couchdbkit

DIR_NAME = os.path.abspath(os.path.dirname(__file__))
VALID_VIEW_SETS = ('gnotify_events',)


def deploy(uri, db_name, view_set):
    server = couchdbkit.Server(uri)

    loader = couchdbkit.FileSystemDocsLoader(
        os.path.join(DIR_NAME, view_set, '_design')
    )

    db = server.get_db(db_name)
    loader.sync(db, verbose=True)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Couch view deployment')
    parser.add_argument(
        '-s', '--server', metavar='uri', type=str, nargs='?', dest='server',
        default='http://localhost:5984/'
    )
    parser.add_argument('-d', '--db', '--database', type=str, dest='db')
    parser.add_argument(
        '-v', '--view-set', type=str, nargs='?', dest='view_set', default=None
    )

    args = parser.parse_args()

    if args.view_set is None:
        args.view_set = args.db

    if args.view_set not in VALID_VIEW_SETS:
        print 'View set %s invalid, pick one of %s' % (
            args.view_set, VALID_VIEW_SETS
        )
        sys.exit(1)

    deploy(args.server, args.db, args.view_set)

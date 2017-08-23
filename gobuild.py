#!/usr/bin/python3
import os
import sys
import argparse
import subprocess
import base64



goosarchs = [
    ('darwin', '386'),
    ('darwin', 'amd64'),
    # # ('darwin', 'arm'),  // not compiling
    # # ('darwin', 'arm64'),  // not compiling
    # ('dragonfly', 'amd64'),
    ('freebsd', '386'),
    ('freebsd', 'amd64'),
    ('freebsd', 'arm'),
    ('linux', '386'),
    ('linux', 'amd64'),
    ('linux', 'arm'),
    ('linux', 'arm64'),
    # ('linux', 'ppc64'),
    # ('linux', 'ppc64le'),
    # ('linux', 'mips'),
    # ('linux', 'mipsle'),
    # ('linux', 'mips64'),
    # ('linux', 'mips64le'),
    # ('netbsd', '386'),
    # ('netbsd', 'amd64'),
    # ('netbsd', 'arm'),
    # ('openbsd', '386'),
    # ('openbsd', 'amd64'),
    # ('openbsd', 'arm'),
    # ('plan9', '386'),
    # ('plan9', 'amd64'),
    # # ('solaris', 'amd64'),  // not compiling
    ('windows', '386'),
    ('windows', 'amd64'),
]


GOFILE = 'main.go'
TARGET = 'siridb-prompt'


def get_version(path):
    version = None
    with open(os.path.join(path, GOFILE), 'r') as f:
        for line in f:
            if line.startswith('const AppVersion ='):
                version = line.split('"')[1]
    if version is None:
        raise Exception('Cannot find version in {}'.format(GOFILE))
    return version


def build_all():
    path = os.path.dirname(__file__)
    version = get_version(path)
    outpath = os.path.join(path, 'bin', version)
    if not os.path.exists(outpath):
        os.makedirs(outpath)

    for goos, goarch in goosarchs:
        tmp_env = os.environ.copy()
        tmp_env["GOOS"] = goos
        tmp_env["GOARCH"] = goarch
        outfile = os.path.join(outpath, '{}_{}_{}_{}.{}'.format(
            TARGET,
            version,
            goos,
            goarch,
            'exe' if goos == 'windows' else 'bin'))
        with subprocess.Popen(
                ['go', 'build', '-o', outfile],
                env=tmp_env,
                cwd=path,
                stdout=subprocess.PIPE) as proc:
            print('Building {}/{}...'.format(goos, goarch))


def build(output=''):
    path = os.path.dirname(__file__)
    version = get_version(path)
    outfile = output if output else os.path.join(path, '{}_{}.{}'.format(
        TARGET, version, 'exe' if sys.platform.startswith('win') else 'bin'))
    args = ['go', 'build', '-o', outfile]

    with subprocess.Popen(
            args,
            cwd=path,
            stdout=subprocess.PIPE) as proc:
        print('Building {}...'.format(outfile))


def install_packages():
    path = os.path.dirname(__file__)
    with subprocess.Popen(
            ['go', 'get', '-d'],
            cwd=path,
            stdout=subprocess.PIPE) as proc:
        print(
            'Downloading required go packages and dependencies.\n'
            '(be patient, this can take some time)...')


if __name__ == '__main__':

    parser = argparse.ArgumentParser()

    parser.add_argument(
        '-i', '--install-packages',
        action='store_true',
        help='install required go packages including dependencies')

    parser.add_argument(
        '-b', '--build',
        action='store_true',
        help='build binary')

    parser.add_argument(
        '-o', '--output',
        default='',
        help='alternative output filename (requires -b/--build)')

    parser.add_argument(
        '-a', '--build-all',
        action='store_true',
        help='build production binaries for all goos and goarchs')

    args = parser.parse_args()

    if args.output and not args.build:
        print('Cannot use -o/--output without -b/--build')
        sys.exit(1)

    if args.install_packages:
        install_packages()
        print('Finished installing required packages and dependencies!')

    if args.build:
        print('Build binary')
        build(output=args.output)
        print('Finished build!')

    if args.build_all:
        build_all()
        print('Finished building binaries!')

    if not any([
            args.install_packages,
            args.build,
            args.build_all]):
        parser.print_usage()

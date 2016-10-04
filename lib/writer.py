'''writer functions.

writer functions are output functions used by SiriDB Prompt.

:copyright: 2016, Jeroen van der Heijden (Transceptor Technology)
'''
import json
import math
import datetime
import operator


time_precision = None


def yellow(s):
    return '\x1b[33m{}\x1b[0m'.format(s)


def bold(s):
    return '\033[1m{}\x1b[0m'.format(s)


def error(s):
    print(yellow('{}\n'.format(s)))


def fmt_large_num(n):
    s = str(n)
    if len(s) <= 6:
        return s
    else:
        start = -(len(s)) % 3
        return ''.join([c if n == start or n % 3 else '.' + c
                        for n, c in enumerate(s, start=start)])


def fmt_size_bytes(size):
    lookup = 'BKMGTPEZYXWVU'
    if size > 0:
        i = int(min(math.log(size) // math.log(1024), 12))
        n = round(size * 100 / 1024 ** i) / 100
        return '{} {}B'.format(n, lookup[i]) if i else '{} bytes'.format(n)
    return '0 bytes'


def fmt_size_mb(size):
    size *= 1024 ** 2
    return fmt_size_bytes(size)


def fmt_duration(seconds):
    if seconds == 1:
        return '1 second'
    if seconds <= 60 * 2:
        return '{} seconds'.format(seconds)
    if seconds <= 3600 * 2:
        return '{} minutes'.format(round(seconds / 60))
    if seconds <= 86400 * 2:
        return '{} hours'.format(round(seconds / 3600))
    if seconds <= 2592000 * 2:
        return '{} days'.format(round(seconds / 86400))
    if seconds <= 31557600 * 2:
        return '{} months'.format(round(seconds / 2592000))
    return '{} years'.format(round(seconds / 31557600))


def fmt_manhole(s):
    return s if 'not installed' in s else \
        '{0}, example usage: socat - unix-connect:{0}'.format(s)


def fmt_percentage(p):
    return '{} ({}%)'.format(p, int(p * 100))


def fmt_utc_timestamp(t):
    if time_precision is None:
        return str(t)
    t /= _FACTOR[time_precision]
    return datetime.datetime.utcfromtimestamp(t).strftime('%Y-%m-%d %H:%M:%SZ')


def fmt_naive_timestamp(t):
    if time_precision is None:
        return str(t)
    t /= _FACTOR[time_precision]
    return datetime.datetime.fromtimestamp(t).strftime('%Y-%m-%d %H:%M:%S')


_COUNT_MAP = {
    'series': fmt_large_num,
    'series_length': fmt_large_num,
    'pools_series': fmt_large_num,
    'servers_received_points': fmt_large_num,
    'groups_series': fmt_large_num,
    'servers_mem_usage': fmt_size_mb,
    'shards_size': fmt_size_bytes,
}

_FMT_MAP = {
    'received_points': fmt_large_num,
    'series': fmt_large_num,
    'mem_usage': fmt_size_mb,
    'size': fmt_size_bytes,
    'length': fmt_large_num,
    'buffer_size': fmt_size_bytes,
    'uptime': fmt_duration,
    'queries_timeout': fmt_duration,
    'manhole': fmt_manhole,
    'drop_threshold': fmt_percentage,
    'start': fmt_utc_timestamp,
    'end': fmt_utc_timestamp,
    'timestamp': fmt_utc_timestamp
}

_FACTOR = {
    's': 1e0,
    'ms': 1e3,
    'us': 1e6,
    'ns': 1e9
}


def print_table(columns, rows, max_rows=None, fmt_by_column=True, sort=True):
    col_sizes = [len(field) for field in columns]
    for row in rows:
        for n, column in enumerate(row):
            col_sizes[n] = max(len(_FMT_MAP.get(columns[n]
                                                if fmt_by_column
                                                else row[0] if n else None,
                                                str)(column)),
                               col_sizes[n])
            if max_rows and n == max_rows:
                break
    print(' \u2503 '.join(['{:<{}}'.format(column, size)
                           for column, size in zip(columns, col_sizes)]))
    print('\u2501\u2547\u2501'.join(['\u2501' * n for n in col_sizes]))
    for n, row in enumerate(sorted(rows, key=operator.itemgetter(0))
                            if sort else
                            rows):
        print(' \u2502 '.join(
            ['{:<{}}'.format(_FMT_MAP.get(columns[i]
                                          if fmt_by_column
                                          else row[0] if i else None,
                                          str)(column), size)
             for i, (column, size) in enumerate(zip(row, col_sizes))]))
        if max_rows and n == max_rows:
            print('...\n(use set_output_raw or set_output_json for the full '
                  'result)')
            break
    print('')


def output_list(result):
    columns = result['columns']
    for opt in ('servers', 'series', 'shards', 'groups', 'pools', 'users',
                'networks'):
        rows = result.get(opt)
        if rows is not None:
            break
    print_table(columns, rows)


def output_show(result):
    columns = ('name', 'value')
    rows = [(r['name'], r['value']) for r in result['data']]
    print_table(columns, rows, fmt_by_column=False)


def output_timeit(result):
    timeit = result['__timeit__']
    columns = ('server', 'time')
    rows = [(t['server'], '{:.6f} seconds'.format(t['time']))
            for t in timeit]
    print_table(columns, rows, sort=False)


def output_points(result):
    columns = ('timestamp', 'value')
    for n, (series_name, points) in enumerate(result.items()):
        if series_name.startswith('__'):
            continue
        print('\"{}\"'.format(bold(series_name)))
        print_table(columns, points, max_rows=100, sort=False)
        if n == 100:
            print('...\n(use set_output_raw or set_output_json for the full '
                  'result)')
            break


def output_calc(result):
    t = result['calc']
    print(t)
    if time_precision:
        print('UTC time: {}'.format(fmt_utc_timestamp(t)))
        print('Local time: {}'.format(fmt_naive_timestamp(t)))
    print('')


def output_connections(cluster):
    columns = ('host', 'port', 'connected', 'available')
    rows = [(conn.host,
             conn.port,
             conn._protocol and conn._protocol._connected,
             conn._protocol and conn._protocol._is_available)
            for conn in cluster._connections]
    print_table(columns, rows)


def output_json(result):
    print(json.dumps(result,
                     sort_keys=True,
                     indent=4,
                     separators=(',', ': ')))


def output_pretty(result):
    if 'columns' in result and \
            result['columns'] and \
            isinstance(result['columns'], (tuple, list)) and \
            isinstance(result['columns'][0], str):
        output_list(result)
    elif 'data' in result and \
            result['data'] and \
            isinstance(result['data'], (tuple, list)) and \
            isinstance(result['data'][0], dict):
        output_show(result)
    elif 'success_msg' in result and \
            isinstance(result['success_msg'], str):
        print('{}\n'.format(result['success_msg']))
    elif 'calc' in result and \
            isinstance(result['calc'], int):
        output_calc(result)
    elif 'help' in result and \
            isinstance(result['help'], str):
        print('{}\n'.format(result['help']))
    elif 'motd' in result and \
            isinstance(result['motd'], str):
        print('{}\n'.format(result['motd']))
    #
    #  Count statements
    #
    else:
        for name in ('shards', 'series', 'servers', 'groups', 'pools', 'users',
                     'servers_received_points', 'shards_size', 'series_length'):
            if name in result and \
                    isinstance(result[name], int):
                print('{}: {}\n'.format(
                    name.replace('_', ' ').capitalize(),
                    _COUNT_MAP.get(name, str)(result[name])))
                break
        else:
            if result:
                output_points(result)

    if '__timeit__' in result:
        output_timeit(result)

'''SiriDB Prompt Loop

This is the main loop used by SiriDB Prompt. As soon as the arguments are
accepted you will get stuck in this loop until you quit.

:copyright: 2016, Jeroen van der Heijden (Transceptor Technology)
'''
import os
import sys
import logging
import asyncio
import json
from prompt_toolkit.shortcuts import prompt_async
from prompt_toolkit.history import InMemoryHistory
from prompt_toolkit.history import FileHistory
from . import writer
from .history import clear_history
from .history import crop_history
from .history import get_history
from .completer import SiriCompleter
from .completer import confirm_completer
from .manager import manager
from .version import __version__
from .version import __version_info__

sys.path.append('/home/joente/workspace/siridb-connector/')
from siridb.connector.lib.exceptions import QueryError
from siridb.connector.lib.exceptions import InsertError

class csvhandler:
    def loads(self, *args, **kwargs):
        raise NotImplementedError('Not yet implemented')

__DEFAULT_OUTPUT = 'PRETTY'
__IMPORT_MAP = {
    'import_json': json,
    'import_csv': csvhandler}
force_exit = set()


def set_loglevel(inp):
    for l in SiriCompleter.LOG_LEVEL_SET:
        if inp == l:
            logger = logging.getLogger()
            level = l.split('_')[-1].upper()
            logger.setLevel(level)
            print('Loglevel set to {!r}\n'.format(level))
            return True
    return False


def set_output(inp):
    for l in SiriCompleter.OUTPUT_SET:
        if inp == l:
            outp = l.split('_')[-1].upper()
            print('Output set to {!r}\n'.format(outp))
            return outp
    return False


async def import_from_file(inp):
    for cmd in SiriCompleter.IMPORTS:
        if inp.startswith(cmd):
            break
    else:
        return None
    fn = inp[len(cmd):].strip()
    if fn:
        c = fn[0]
        if c in '\'"':
            fn = fn[1:-1].replace(c * 2, c)
    if not fn:
        raise FileNotFoundError('Expecting a filename, for example: '
                                '\'{} ./my_file.{}\''
                                .format(cmd, cmd.split('_')[-1]))
    with open(fn, 'r', encoding='utf-8') as f:
        data = __IMPORT_MAP[cmd].load(f)

    while True:
        answer = await prompt_async('Are you sure you want to insert data for '
                                    '{} series? (yes/no)\n> '
                                    .format(len(data)),
                                    patch_stdout=True,
                                    key_bindings_registry=manager.registry,
                                    completer=confirm_completer)
        if answer in ('n', 'no'):
            raise ValueError('Import cancelled by user.')
        if answer in ('y', 'yes'):
            break
        print(writer.yellow('Expecting \'yes\' or \'no\' but got {!r}'
                            .format(answer)))

    return data


async def stop(cluster):
    cluster.close()
    pending = {task
               for task in asyncio.Task.all_tasks()
               if task is not asyncio.Task.current_task()}
    if pending:
        _, pending = await asyncio.wait(pending,
                                        timeout=3)


async def prompt_loop(cluster):
    output = __DEFAULT_OUTPUT
    connections = await cluster.connect()

    if None not in connections:
        logging.error('Could not setup a connection')
        await stop(cluster)
        return

    completer = SiriCompleter()

    try:
        path = os.path.join(os.path.expanduser('~'), '.siridb-prompt')
        if not os.path.isdir(path):
            os.mkdir(path)

        history_file = os.path.join(path,
                                    '{}@{}.history'.format(cluster._username,
                                                           cluster._dbname))
        history = FileHistory(history_file)
    except Exception as e:
        history = InMemoryHistory()

    try:
        result = await cluster.query('show time_precision, version', timeout=5)
        writer.time_precision = result['data'][0]['value']
        version = result['data'][1]['value']
        if tuple(map(int, version.split('.')[:2])) != __version_info__[:2]:
            logging.warning('SiriDB Prompt is version is version {} while at '
                            'least one server your connecting to is running '
                            'version {}. This is not necessary a problem but '
                            'we recommend using the same version.'
                            .format(__version__, version))
    except Exception as e:
        logging.warning('Could not read the version and time precision, '
                        'timestamps will not be formatted and we are not sure '
                        'if the database is running a correct version. '
                        '(reason: {})'.format(e))

    sys.stdout.write('\x1b]2;{}@{}\x07'.format(cluster._username,
                                               cluster._dbname))

    while True:
        force_exit.clear()
        method = cluster.query
        inp = await prompt_async('{}@{}> '.format(cluster._username,
                                                  cluster._dbname),
                                 patch_stdout=True,
                                 key_bindings_registry=manager.registry,
                                 history=history,
                                 completer=completer)
        inp = inp.strip()

        if inp == 'exit':
            break

        if get_history(inp, history):
            continue

        if inp == 'clear':
            clear_history(history)
            print('History is cleared.\n')
            continue

        if inp == 'connections':
            writer.output_connections(cluster)
            continue

        if inp == 'get_output':
            print('{}\n'.format(output))
            continue

        if inp == 'get_loglevel':
            print('{}\n'.format(
                logging._levelToName[logging.getLogger().level]))
            continue

        if set_loglevel(inp):
            continue

        o = set_output(inp)
        if o:
            output = o
            continue

        try:
            data = await import_from_file(inp)
        except Exception as e:
            writer.error('Import error: {}'.format(e))
            continue
        else:
            if data is not None:
                inp = data
                method = cluster.insert

        try:
            result = await method(inp)
        except (QueryError, InsertError) as e:
            writer.error(e)
        except TimeoutError as e:
            logging.critical(
                '{} (You might receive another error \'Package ID '
                'not found\' if the result arrives)'.format(e))
        except asyncio.CancelledError:
            logging.critical('Request is cancelled...')
        except Exception as e:
            if str(e):
                logging.critical(e)
            else:
                logging.exception(e)
        else:
            if isinstance(result, str) or output == 'RAW':
                print('{}\n'.format(result))
            elif output == 'JSON':
                writer.output_json(result)
            elif output == 'PRETTY':
                writer.output_pretty(result)
            else:
                raise TypeError('Unknown output set: {}'.format(output))

    await stop(cluster)
    crop_history(history)

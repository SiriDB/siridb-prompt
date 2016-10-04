'''History for SiriDB Prompt.

This are methods used to 'control' the history used by SiriDB Prompt.

:copyright: 2016, Jeroen van der Heijden (Transceptor Technology)
'''
import logging
import re
import functools
from lib import writer


__MAX_ITEMS_IN_HISTORY = 3000


def clear_history(history):
    if hasattr(history, 'filename'):
        with open(history.filename, 'wb'):
            pass
    del history.strings[:]


def crop_history(history, ma=__MAX_ITEMS_IN_HISTORY):
    if not hasattr(history, 'filename'):
        return
    n = len(history.strings)
    if n > ma:
        logging.info('Got {} items in history, we will crop history and '
                     'save the last {} lines ...'.format(n, ma))
        content = []
        i = 0
        got_line = False
        with open(history.filename, 'rb') as f:
            for line in f:
                if i >= n - ma:
                    content.append(line)
                line = line.decode('utf-8')
                if line.startswith('+'):
                    got_line = True
                else:
                    if got_line:
                        i += 1
                    got_line = False

        with open(history.filename, 'wb') as f:
            f.write(b''.join(content))


def _extract_regex(r):
    def wrap(s):
        return '^{}$'.format(s)
    return \
        (wrap(r[1:-2]), re.IGNORECASE) if r[-1] == 'i' else (wrap(r[1:-1]), 0)


def get_history(inp, history, _re_match=re.compile('^/.+/i?$')):
    if not inp.startswith('history'):
        return

    inp = inp[len('history'):].strip()
    if not inp:
        for item in history:
            print(item)
        print('\nGot {} item(s) in history, type \'clear\' to erase them\n'
              .format(len(history)))
        return True

    if not re.match(_re_match, inp):
        writer.error('Invalid regular expression: {!r}, expecting something '
                     'like /select.*/'.format(inp))
        return True

    try:
        test = functools.partial(re.match, re.compile(*_extract_regex(inp)))
    except re.error:
        writer.error('Error found in regular expression: {!r}'.format(inp))
        return True

    matches = tuple(filter(test, history))

    for item in matches:
        print(item)
    print('\nMatches {} out of {} item(s) in history.\n'
          .format(len(matches), len(history)))
    return True

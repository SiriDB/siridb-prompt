'''Export python grammar to C grammar files

Author: Jeroen van der Heijden (Transceptor Technology)
Date: 2016-10-10
'''
import os
from lib.grammar import siri_grammar

if __name__ == '__main__':
    c_file, h_file = siri_grammar.export_c(target='siri/grammar/grammar')

    EXPOTR_PATH = 'cgrammar'

    try:
        os.makedirs(EXPOTR_PATH)
    except FileExistsError:
        pass

    with open(os.path.join(EXPOTR_PATH, 'grammar.c'),
            'w',
            encoding='utf-8') as f:
        f.write(c_file)
    with open(os.path.join(EXPOTR_PATH, 'grammar.h'),
            'w',
            encoding='utf-8') as f:
        f.write(h_file)

    print('\nFinished creating new c-grammar files...\n')

    js_file = siri_grammar.export_js()

    EXPOTR_PATH = 'jsgrammar'

    try:
        os.makedirs(EXPOTR_PATH)
    except FileExistsError:
        pass

    with open(os.path.join(EXPOTR_PATH, 'grammar.js'),
            'w',
            encoding='utf-8') as f:
        f.write(js_file)

    print('\nFinished creating new js-grammar file...\n')
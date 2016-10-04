'''Completer for SiriDB Prompt.

SiriCompleter class used for auto-completion in the prompt.
We use the Siri Grammar and some additional prompt commands for full auto-
completion support.

:copyright: 2016, Jeroen van der Heijden (Transceptor Technology)
'''
from prompt_toolkit.completion import Completer
from prompt_toolkit.completion import Completion
from prompt_toolkit.contrib.completers import WordCompleter
from .grammar import siri_grammar
from pyleri import Keyword


confirm_completer = WordCompleter(['yes', 'no'])


class SiriCompleter(Completer):

    LOG_LEVEL_SET = (
        'set_loglevel_debug',
        'set_loglevel_info',
        'set_loglevel_warning',
        'set_loglevel_error',
        'set_loglevel_critical')

    OUTPUT_SET = (
             'set_output_json',
             'set_output_pretty',
             'set_output_raw')

    IMPORTS = (
        'import_csv',
        'import_json')

    WORDS = ('exit',
             'get_output',
             'get_loglevel',
             'history',
             'connections',
             'clear') + LOG_LEVEL_SET + OUTPUT_SET + IMPORTS

    def get_completions(self, document, complete_event):
        query = document.text_before_cursor
        for word in self.WORDS:
            if word.startswith(query):
                yield Completion(word, start_position=-len(query))

        result = siri_grammar.parse(query)
        rest = query[result.pos:]
        for elem in result.expecting:
            if isinstance(elem, Keyword):
                word = str(elem)
                if not rest and query and query[-1].isspace() or \
                        rest and word.startswith(rest):
                    yield Completion(
                        word + ' ',
                        display=word,
                        start_position=-len(rest))

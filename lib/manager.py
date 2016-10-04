'''Key manager

Used to bind CTRL-C to the prompt and prevent quiting the application.

:copyright: 2016, Jeroen van der Heijden (Transceptor Technology)
'''
import logging
from prompt_toolkit.key_binding.manager import KeyBindingManager
from prompt_toolkit.keys import Keys


manager = KeyBindingManager.for_prompt()


@manager.registry.add_binding(Keys.ControlC)
def _(event):
    def press_ctrl_c():
        logging.warning('You pressed CTRL+C, type \'exit\' if you want to '
                        'quit')
    event.cli.run_in_terminal(press_ctrl_c)

'''
Setup siridb-prompt (main.py) using cx_Freeze

Installation cx_Freeze on Ubuntu:

see: https://bitbucket.org/anthony_tuininga/cx_freeze/
        issue/32/cant-compile-cx_freeze-in-ubuntu-1304

 - sudo apt-get install python3-dev
 - sudo apt-get install libssl-dev
 - Open setup.py and change the line
     if not vars.get("Py_ENABLE_SHARED", 0):
   to
     if True:
 - python3 setup.py build
 - sudo python3 setup.py install

'''
import os
import platform
from cx_Freeze import setup, Executable
from lib.version import __version__

architecture = {'64bit': 'x86_64', '32bit': 'i386'}[platform.architecture()[0]]

ccsv_fn = 'ccsv-{0}.cpython-35m-{0}-linux-gnu.so'.format(
    {'64bit': 'x86_64', '32bit': 'i386'}[platform.architecture()[0]])

build_exe_options = {
    'build_exe': os.path.join('build', str(__version__)),
    'packages': [
        'encodings',
        'os',
        'argparse',
        'signal',
        'functools',
        'passlib',
        'asyncio',
        'shutil',
        'logging',
        'getpass',
        'collections',
        'uuid',
        'pickle',
        'pyleri',
        'prompt_toolkit',
        'pygments',
        'json'],
    'excludes': [
        'django',
        'google',
        'twisted'],
    'optimize': 2,
    'include_files': [
        ('help', 'help'),
        (os.path.join('ccsv', ccsv_fn), ccsv_fn)]}


setup(
    name='prompt',
    version=__version__,
    description='Prompt for SiriDB',
    options={'build_exe': build_exe_options},
    executables=[Executable('main.py', targetName='siridb-prompt')])

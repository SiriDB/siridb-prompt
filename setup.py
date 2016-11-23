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
from cx_Freeze import setup, Executable
from lib.version import __version__


build_exe_options = {
    'build_exe': os.path.join('build', __version__),
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
        'json',
        'csvloader',
        'qpack',
        'siridb-connector'],
    'excludes': [
        'django',
        'google',
        'twisted'],
    'optimize': 2,
    'include_files': []}


setup(
    name='siridb-prompt',
    version=__version__,
    description='Prompt for SiriDB',
    options={'build_exe': build_exe_options},
    executables=[Executable('main.py', targetName='siridb-prompt')])

from setuptools import setup, find_packages

setup(
    name='jbot',
    version='0.0.1',
    description='Jobs Bot',
    packages=find_packages(),
    python_requires='>=3.8, <4',
    package_data={
        'jbot.log': [
            'logger.yaml',
        ],
    },
    install_requires=[
        'click',
        'pyyaml',
        'python-dotenv',

        'aiogram',

        'grpcio',
        'grpcio-tools',
    ],
    entry_points={
        'console_scripts': [
            'jb=jbot.cli:cli'
        ],
    },
)

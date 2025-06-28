from setuptools import setup, find_packages

setup(
    name="analytics-service",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "flask>=2.0.0",
        "pandas>=1.3.0",
        "numpy>=1.20.0",
        "matplotlib>=3.4.0",
        "requests>=2.25.0",
    ],
    python_requires=">=3.8",
)

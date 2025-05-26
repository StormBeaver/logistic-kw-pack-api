import setuptools

setuptools.setup(
    name="grpc-logistic-pack-api",
    version="1.0.0",
    author="rusdevop",
    author_email="rusdevops@gmail.com",
    description="GRPC python client for logistic-pack-api",
    url="https://github.com/ozonmp/logistic-pack-api",
    packages=setuptools.find_packages(),
    package_data={"ozonmp.logistic_pack_api.v1": ["logistic_pack_api_pb2.pyi"]},
    python_requires='>=3.5',
)
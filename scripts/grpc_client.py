import asyncio

from grpclib.client import Channel

from ozonmp.logistic_pack_api.v1.logistic_pack_api_grpc import LogisticPackApiServiceStub
from ozonmp.logistic_pack_api.v1.logistic_pack_api_pb2 import DescribePackV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = LogisticPackApiServiceStub(channel)

        req = DescribePackV1Request(pack_id=1)
        reply = await client.DescribePackV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())

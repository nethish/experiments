# -*- coding: utf-8 -*-

from pprint import pprint
from time import sleep

from pymongo import MongoClient
from pymongo.errors import PyMongoError


# Optional Resume After
# The resume after only works as long as the the operation is present in the oplog
# The data is streamed from oplog, and oplog is a feature of a replica set. So you need to configure mongo in replica mode for it to work
resume_after = {
    "_data": "8267FEC222000000022B022C0100296E5A100483C25F90B0B94EDA808A86AD9B08FC7546645F6964006467FEC22279DA330F90A00AA90004"
}


def main():
    client = MongoClient(
        host="localhost",
        port=27017,
        directConnection=True,
        username=None,
        password=None,
    )

    while True:
        try:
            print("Connecting to MongoDB server and waiting for streams events...")

            # You can replace resume after here to resume the data from last checkpoint
            with client.watch(pipeline=None, resume_after=None) as stream:
                for event in stream:
                    print(
                        f"Received event: {event.get('operationType')} "
                        f"for document: {event.get('documentKey', {}).get('_id')}."
                    )

                    pprint(event)

        except PyMongoError as error:
            print(
                f"An error has been raised. "
                f"Trying to resume after 30 seconds. "
                f"Error: {error}."
            )

            sleep(30)


if __name__ == "__main__":
    main()

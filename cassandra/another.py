from cassandra.cluster import Cluster
from cassandra.query import SimpleStatement
from cassandra import ConsistencyLevel

# Connect to Cassandra cluster nodes on different ports
# Note: Cassandra nodes usually listen on 9042 by default,
# but to connect on different ports, you specify each host:port tuple.
# The driver expects a list of contact points (hostnames/IPs).
# To specify ports explicitly, we need to set up a custom connection for each node.
# However, Cluster API accepts only hosts (no port).
# So we will create separate Cluster connections per port and coordinate manually,
# or better, use a single port where nodes expose the native protocol (usually 9042).

# --- For demonstration, if your cluster nodes use different ports for native protocol ---
# We create four separate Cluster connections for each node.

# But it's unusual: normally Cassandra cluster nodes use same port (9042).
# If you really want to connect to different ports on localhost,
# you'll have to create separate Cluster sessions.


def create_session(host, port=9042):
    return Cluster([host], port=port).connect()


# Create sessions per node
sessions = [
    create_session("127.0.0.1", 9042),
    create_session("127.0.0.1", 9043),
    create_session("127.0.0.1", 9044),
    create_session("127.0.0.1", 9045),
]

# We'll just use session 0 as coordinator for simplicity
session = sessions[0]

KEYSPACE = "testks"
TABLE = "users"


def setup_keyspace_and_table():
    # Create keyspace with RF=3 for example
    session.execute(f"""
        CREATE KEYSPACE IF NOT EXISTS {KEYSPACE}
        WITH replication = {{ 'class': 'SimpleStrategy', 'replication_factor': '3' }}
    """)
    session.set_keyspace(KEYSPACE)

    session.execute(f"""
        CREATE TABLE IF NOT EXISTS {TABLE} (
            user_id text PRIMARY KEY,
            name text,
            email text
        )
    """)


def insert_user(user_id, name, email):
    query = f"INSERT INTO {TABLE} (user_id, name, email) VALUES (%s, %s, %s)"
    session.execute(query, (user_id, name, email))
    print(f"Inserted user {user_id}")


def select_user(user_id):
    query = f"SELECT * FROM {TABLE} WHERE user_id=%s"
    row = session.execute(query, (user_id,)).one()
    if row:
        print(f"User found: id={row.user_id}, name={row.name}, email={row.email}")
        return row
    else:
        print("User not found")
        return None


def update_user_email(user_id, new_email):
    query = f"UPDATE {TABLE} SET email=%s WHERE user_id=%s"
    session.execute(query, (new_email, user_id))
    print(f"Updated user {user_id} email to {new_email}")


def delete_user(user_id):
    query = f"DELETE FROM {TABLE} WHERE user_id=%s"
    session.execute(query, (user_id,))
    print(f"Deleted user {user_id}")


def main():
    setup_keyspace_and_table()

    # Create
    insert_user("user123", "Alice", "alice@example.com")

    # Read
    select_user("user123")

    # Update
    update_user_email("user123", "alice_new@example.com")

    # Read again to confirm update
    select_user("user123")

    # Delete
    delete_user("user123")

    # Confirm deletion
    select_user("user123")


if __name__ == "__main__":
    main()

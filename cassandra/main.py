#!/usr/bin/env python3
"""
Cassandra Python Client
Demonstrates inserting and retrieving data from Cassandra cluster
"""

from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider
from cassandra.policies import DCAwareRoundRobinPolicy
import uuid
import json
from datetime import datetime
import time


class CassandraClient:
    def __init__(self):
        # Connection configuration
        self.hosts = ["localhost"]  # Docker host
        self.ports = [9042, 9043, 9044, 9045]  # All node ports
        self.keyspace = "demo_keyspace"
        self.cluster = None
        self.session = None

    def connect(self):
        """Connect to Cassandra cluster"""
        try:
            # Since your test script works with port 9042, let's use that approach
            print("Attempting to connect to Cassandra cluster on port 9042...")

            self.cluster = Cluster(
                contact_points=["127.0.0.1"],
                port=9042,
                connect_timeout=10,
                control_connection_timeout=10,
            )
            self.session = self.cluster.connect()
            print("✅ Connected to Cassandra cluster successfully!")
            return True

        except Exception as e:
            print(f"❌ Connection failed: {e}")
            print("\n🔍 Troubleshooting tips:")
            print("1. Check if containers are running: docker ps")
            print("2. Check cluster status: docker exec cassandra1 nodetool status")
            print("3. Wait for cluster to be ready (can take 2-5 minutes)")
            print("4. Check logs: docker logs cassandra1")
            return False

    def setup_keyspace_and_tables(self):
        """Create keyspace and tables"""
        try:
            # Create keyspace
            create_keyspace = f"""
            CREATE KEYSPACE IF NOT EXISTS {self.keyspace}
            WITH REPLICATION = {{
                'class': 'NetworkTopologyStrategy',
                'datacenter1': 3
            }}
            """
            self.session.execute(create_keyspace)
            print(f"✅ Keyspace '{self.keyspace}' created/verified")

            # Use the keyspace
            self.session.set_keyspace(self.keyspace)

            # Create users table
            create_users_table = """
            CREATE TABLE IF NOT EXISTS users (
                user_id UUID PRIMARY KEY,
                username TEXT,
                email TEXT,
                age INT,
                created_at TIMESTAMP,
                metadata MAP<TEXT, TEXT>
            )
            """
            self.session.execute(create_users_table)
            print("✅ Users table created/verified")

            # Create orders table (demonstrates composite key)
            create_orders_table = """
            CREATE TABLE IF NOT EXISTS orders (
                user_id UUID,
                order_id UUID,
                product_name TEXT,
                price DECIMAL,
                order_date TIMESTAMP,
                PRIMARY KEY (user_id, order_id)
            ) WITH CLUSTERING ORDER BY (order_id DESC)
            """
            self.session.execute(create_orders_table)
            print("✅ Orders table created/verified")

        except Exception as e:
            print(f"Setup error: {e}")
            raise

    def insert_sample_data(self):
        """Insert sample data into tables"""
        try:
            # Sample users data
            users_data = [
                {
                    "user_id": uuid.uuid4(),
                    "username": "john_doe",
                    "email": "john@example.com",
                    "age": 30,
                    "metadata": {"country": "USA", "city": "New York"},
                },
                {
                    "user_id": uuid.uuid4(),
                    "username": "jane_smith",
                    "email": "jane@example.com",
                    "age": 25,
                    "metadata": {"country": "Canada", "city": "Toronto"},
                },
                {
                    "user_id": uuid.uuid4(),
                    "username": "bob_wilson",
                    "email": "bob@example.com",
                    "age": 35,
                    "metadata": {"country": "UK", "city": "London"},
                },
            ]

            # Insert users
            insert_user_stmt = self.session.prepare("""
                INSERT INTO users (user_id, username, email, age, created_at, metadata)
                VALUES (?, ?, ?, ?, ?, ?)
            """)

            for user in users_data:
                self.session.execute(
                    insert_user_stmt,
                    [
                        user["user_id"],
                        user["username"],
                        user["email"],
                        user["age"],
                        datetime.now(),
                        user["metadata"],
                    ],
                )

            print(f"✅ Inserted {len(users_data)} users")

            # Insert orders for first user
            user_id = users_data[0]["user_id"]
            orders_data = [
                {
                    "order_id": uuid.uuid4(),
                    "product_name": "Laptop",
                    "price": 999.99,
                },
                {
                    "order_id": uuid.uuid4(),
                    "product_name": "Mouse",
                    "price": 29.99,
                },
                {
                    "order_id": uuid.uuid4(),
                    "product_name": "Keyboard",
                    "price": 79.99,
                },
            ]

            insert_order_stmt = self.session.prepare("""
                INSERT INTO orders (user_id, order_id, product_name, price, order_date)
                VALUES (?, ?, ?, ?, ?)
            """)

            for order in orders_data:
                self.session.execute(
                    insert_order_stmt,
                    [
                        user_id,
                        order["order_id"],
                        order["product_name"],
                        order["price"],
                        datetime.now(),
                    ],
                )

            print(
                f"✅ Inserted {len(orders_data)} orders for user {users_data[0]['username']}"
            )

            return users_data[0]["user_id"]  # Return first user ID for queries

        except Exception as e:
            print(f"Insert error: {e}")
            raise

    def retrieve_data(self, sample_user_id):
        """Retrieve and display data"""
        try:
            print("\n" + "=" * 50)
            print("RETRIEVING DATA")
            print("=" * 50)

            # 1. Get all users
            print("\n📋 All Users:")
            users_result = self.session.execute("SELECT * FROM users")
            for user in users_result:
                print(f"  • {user.username} ({user.email}) - Age: {user.age}")
                print(f"    ID: {user.user_id}")
                print(f"    Created: {user.created_at}")
                print(
                    f"    Metadata: {dict(user.metadata) if user.metadata else 'None'}"
                )
                print()

            # 2. Get specific user by ID
            print(f"\n🔍 User Details (ID: {sample_user_id}):")
            user_stmt = self.session.prepare("SELECT * FROM users WHERE user_id = ?")
            user_result = self.session.execute(user_stmt, [sample_user_id])
            user = user_result.one()
            if user:
                print(f"  Username: {user.username}")
                print(f"  Email: {user.email}")
                print(f"  Age: {user.age}")
                print(f"  Metadata: {dict(user.metadata) if user.metadata else 'None'}")

            # 3. Get orders for specific user
            print(f"\n🛒 Orders for User {user.username}:")
            orders_stmt = self.session.prepare("SELECT * FROM orders WHERE user_id = ?")
            orders_result = self.session.execute(orders_stmt, [sample_user_id])
            total_amount = 0
            for order in orders_result:
                print(f"  • {order.product_name}: ${order.price}")
                print(f"    Order ID: {order.order_id}")
                print(f"    Date: {order.order_date}")
                total_amount += float(order.price)
                print()
            print(f"  💰 Total Orders Value: ${total_amount:.2f}")

            # 4. Count operations
            print("\n📊 Statistics:")
            users_count = self.session.execute("SELECT COUNT(*) FROM users").one().count
            orders_count = (
                self.session.execute("SELECT COUNT(*) FROM orders").one().count
            )
            print(f"  Total Users: {users_count}")
            print(f"  Total Orders: {orders_count}")

        except Exception as e:
            print(f"Retrieve error: {e}")
            raise

    def demonstrate_updates_deletes(self, sample_user_id):
        """Demonstrate UPDATE and DELETE operations"""
        try:
            print("\n" + "=" * 50)
            print("UPDATE & DELETE OPERATIONS")
            print("=" * 50)

            # Update user age
            print("\n🔄 Updating user age...")
            update_stmt = self.session.prepare(
                "UPDATE users SET age = ? WHERE user_id = ?"
            )
            self.session.execute(update_stmt, [31, sample_user_id])

            # Verify update
            user_result = self.session.execute(
                self.session.prepare(
                    "SELECT username, age FROM users WHERE user_id = ?"
                ),
                [sample_user_id],
            )
            user = user_result.one()
            print(f"  ✅ Updated {user.username}'s age to {user.age}")

            # Add new order
            print("\n➕ Adding new order...")
            new_order_stmt = self.session.prepare("""
                INSERT INTO orders (user_id, order_id, product_name, price, order_date)
                VALUES (?, ?, ?, ?, ?)
            """)
            new_order_id = uuid.uuid4()
            self.session.execute(
                new_order_stmt,
                [sample_user_id, new_order_id, "Headphones", 149.99, datetime.now()],
            )
            print("  ✅ Added new order: Headphones ($149.99)")

            # Delete the new order
            print("\n🗑️  Deleting the new order...")
            delete_stmt = self.session.prepare(
                "DELETE FROM orders WHERE user_id = ? AND order_id = ?"
            )
            self.session.execute(delete_stmt, [sample_user_id, new_order_id])
            print("  ✅ Deleted Headphones order")

        except Exception as e:
            print(f"Update/Delete error: {e}")
            raise

    def show_cluster_info(self):
        """Display cluster information"""
        try:
            print("\n" + "=" * 50)
            print("CLUSTER INFORMATION")
            print("=" * 50)

            # Get cluster name and nodes
            cluster_info = self.session.execute("SELECT * FROM system.local").one()
            print(f"\n🏢 Cluster Name: {cluster_info.cluster_name}")
            print(f"📍 Data Center: {cluster_info.data_center}")
            print(f"🗄️  Rack: {cluster_info.rack}")

            # Get peer information
            peers = self.session.execute(
                "SELECT peer, data_center, rack FROM system.peers"
            )
            print(f"\n👥 Peer Nodes:")
            for peer in peers:
                print(f"  • {peer.peer} (DC: {peer.data_center}, Rack: {peer.rack})")

        except Exception as e:
            print(f"Cluster info error: {e}")

    def close(self):
        """Close connections"""
        if self.cluster:
            self.cluster.shutdown()
            print("\n🔌 Connection closed")


def main():
    """Main execution function"""
    client = CassandraClient()

    try:
        # Connect to cluster
        if not client.connect():
            return

        # Setup keyspace and tables
        client.setup_keyspace_and_tables()

        # Insert sample data
        sample_user_id = client.insert_sample_data()

        # Wait a moment for data to propagate
        print("\n⏳ Waiting for data to propagate across cluster...")
        time.sleep(2)

        # Retrieve and display data
        client.retrieve_data(sample_user_id)

        # Demonstrate updates and deletes
        client.demonstrate_updates_deletes(sample_user_id)

        # Show cluster information
        client.show_cluster_info()

        print("\n🎉 Demo completed successfully!")

    except Exception as e:
        print(f"❌ Error during execution: {e}")
    finally:
        # Always close connection
        client.close()


if __name__ == "__main__":
    print("🚀 Starting Cassandra Python Demo")
    print("Make sure your Cassandra cluster is running!")
    print("Install required package: pip install cassandra-driver")
    print("-" * 50)
    main()

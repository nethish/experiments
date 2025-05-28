#!/usr/bin/env python3
"""
Simple Cassandra Connection Test
Tests basic connectivity to Cassandra cluster
"""

from cassandra.cluster import Cluster
import time

def test_connection():
    """Test connection to Cassandra"""
    ports_to_test = [9042, 9043, 9044, 9045]
    
    print("🔍 Testing Cassandra connections...")
    print("Make sure Docker containers are running first!")
    print("-" * 50)
    
    for port in ports_to_test:
        try:
            print(f"Testing connection to localhost:{port}...")
            
            cluster = Cluster(
                contact_points=['127.0.0.1'],
                port=port,
                connect_timeout=5
            )
            
            session = cluster.connect()
            
            # Test a simple query
            rows = session.execute("SELECT release_version FROM system.local")
            version = rows.one().release_version
            
            print(f"✅ SUCCESS - Connected to localhost:{port}")
            print(f"   Cassandra version: {version}")
            
            cluster.shutdown()
            return True
            
        except Exception as e:
            print(f"❌ FAILED - localhost:{port}: {e}")
            continue
    
    print("\n❌ Could not connect to any Cassandra node")
    print("\n🔧 Troubleshooting steps:")
    print("1. Check containers: docker ps")
    print("2. Check status: docker exec cassandra1 nodetool status")
    print("3. Check logs: docker logs cassandra1")
    print("4. Wait longer - Cassandra takes time to start")
    return False

if __name__ == "__main__":
    test_connection()

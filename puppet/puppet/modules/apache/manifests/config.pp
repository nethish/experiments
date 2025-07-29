class apache::config {
  file { '/var/www/html/index.html':
    ensure  => file,
    # content => 'Hello from apache::config!',
    source => 'puppet:///modules/apache/index.html',
  }
}

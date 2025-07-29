class apache::config {
  $port      = lookup('apache::port')

  file { '/var/www/html/index.html':
    ensure  => file,
    source => 'puppet:///modules/apache/index.html',
  }

  file { '/etc/apache2/sites-available/000-default.conf':
    ensure  => file,
    content => template('apache/vhost.conf.erb'),
    owner   => 'root',
    group   => 'root',
    mode    => '0644',
    notify  => Service['apache2'],
  }

}

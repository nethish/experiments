class apache {
  include apache::install
  include apache::config
}

# class apache {
#   package { 'apache2':
#     ensure => installed,
#   }
# 
#   service { 'apache2':
#     ensure => running,
#     enable => true,
#   }
# 
#   file { '/var/www/html/index.html':
#     ensure  => file,
#     content => 'Welcome to Apacheeeeeeeeee!',
#     owner   => 'root',
#     group   => 'root',
#     mode    => '0644',
#   }
# }

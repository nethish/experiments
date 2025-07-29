node default {
  file { '/tmp/puppet_was_here':
    ensure  => present,
    content => "Managed by Puppet\n",
    mode    => '0644',
  }

  include apache
}

node 'agent1.puppet' {
  file { '/tmp/hinethish':
    ensure  => present,
    content => "Hi nethish!",
    mode    => '0644',
  }

  file { '/tmp/byenethish':
    ensure  => present,
    content => "bye nethish!",
    mode    => '0644',
  }

  include apache
}


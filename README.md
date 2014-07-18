chef-runner
===========

The purpose of chef-runner explained in one tweet:

> When it comes to dev/testing, fast feedback is everything. Vagrant needs >5s
> before it even starts provisioning. Each time. Let's fix that. -- [@mlafeldt
> (1 Nov 2013)](https://twitter.com/mlafeldt/status/396299646425137152)

For more background, check out my blog post *[Telling people about
chef-runner][blog post]*.

## What is chef-runner?

* A command-line tool that speeds up your Chef development and testing workflow.
* A fast alternative to the painfully slow `vagrant provision` ([demo video]
  comparing both tools).
* [Integrates with Vim](#use-with-vim) so you don't have to leave your editor
  while hacking on recipes.
* Allows you to change infrastructure code and get **immediate feedback**.

## How does it work?

* Prepares Chef configuration and cookbooks in local `.chef-runner` folder.
* Directly executes Chef Solo over SSH (using `vagrant ssh` or OpenSSH).
* Overrides Chef runlist to selectively run recipes.
* Installs cookbook dependencies with Berkshelf and updates changes with rsync.

## Requirements

To use chef-runner, you need the following software:

* [Go] to download, compile, and install the chef-runner tool
* [VirtualBox] or whatever you use with Vagrant
* [Vagrant] - version 1.3.4 or higher
* [Berkshelf] - installable via `gem install berkshelf`
* `rsync`
* `ssh`

Additionally, your cookbook must have the following files:

* `metadata.rb` - must define the cookbook's name
* `Berksfile` - must contain the `metadata` source

To give you an example, the [Practicing Ruby cookbook][pr-cookbook] is known to
work well with chef-runner.

## Installation

First, make sure you have [Go] installed.

To install the `chef-runner` command-line tool run this:

    $ go get github.com/mlafeldt/chef-runner/cmd/chef-runner

For this command to work, `$GOPATH` must be set correctly. Also check that
`$GOPATH` is part of `$PATH`, so that the `chef-runner` executable can be found.
For example, here are the relevant lines from my `~/.bashrc` file:

    export GOPATH="$HOME/devel/go"
    export GOROOT="$(go env GOROOT)"
    export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

## Usage

### Command Line Reference

chef-runner is a simple command-line tool that has a couple of options:

    Usage: chef-runner [options] [--] [<recipe>...]

        -h              Show help text
        -H <name>       Set hostname for direct SSH access
        -M <name>       Set name of Vagrant virtual machine

    Options that will be passed to Chef Solo:

        -F <format>     Set output format (null, doc, minimal, min)
                        default: null
        -l <level>      Set log level (debug, info, warn, error, fatal)
                        default: info
        -j <file>       Load attributes from a JSON file

### Running Chef Recipes

Before you use chef-runner, make sure that the Vagrant machine you want to
provision is running in the background. You can check the status with `vagrant
status`. If the machine isn't up yet, run `vagrant up`.

chef-runner executes one or more recipes you pass on the command line, in the
exact order given. The tool has a flexible recipe syntax allowing you to compose
your [run list] in multiple ways.

1) Run default recipe when passing no arguments:

    $ chef-runner

2) Run local recipe when passing filename:

    $ chef-runner recipes/foo.rb

3) Run local recipe when passing recipe name:

    $ chef-runner foo

4) Run external recipe when passing `cookbook::recipe`:

    $ chef-runner dogs::bar

5) Run multiple recipes (of mixed type) in order given:

    $ chef-runner recipes/foo.rb bar dogs::baz

Moreover, chef-runner allows you to load node attributes from a JSON file (that
is located inside your cookbook):

    $ chef-runner -j chef.json

Last but not least, you can configure both output format and log level of Chef:

    $ chef-runner -F doc -l warn

### Setting up SSH

chef-runner uses `ssh` to execute commands on Vagrant machines. If your
`Vagrantfile` only defines a single machine, you don't need to worry about SSH
at all -- simply run `chef-runner` and it should work.

In a multi-machine environment, you need to specify what Vagrant machine you
want to use. There are two ways to do this.

1) Use the `-M` option to set the name of the Vagrant machine. The machine name
is the name you have defined in your `Vagrantfile`. To get a list of all machine
names, run `vagrant status`.

Example:

    $ chef-runner -M db

2) Use the `-H` option to set a hostname that was configured for direct SSH
access to the Vagrant machine. This requires that your machine has a static IP
address and that your `~/.ssh/config` file has a configuration section for that
hostname (`vagrant ssh-config` can help you here). While this option needs more
setup work, SSH access is a bit faster compared to `-M`.

Example:

    $ chef-runner -H example.local

### Use with Vim

As a matter of fact, I primarily wrote chef-runner for use with Vim. Instead of
jumping back and forth between editing a Chef recipe and running `vagrant
provision`, I wanted to be able to change some code and get immediate feedback
without having to leave the editor. This is where chef-runner's ability to run a
single recipe file -- the file currently open in Vim -- comes in handy.

There's no Vim plugin (yet). For now, you can just stick this one-liner in your
`.vimrc`:

```vim
nnoremap <leader>r :w\|!chef-runner %<cr>
```

With this key mapping in place, make sure that the Vagrant machine is up and
open a recipe in Vim:

```sh
$ vagrant up
$ vim recipes/default.rb
```

Now whenever you type `<leader>r` (your [leader key] then `r`), chef-runner will
run the *current* recipe inside the Vagrant machine, giving you fast feedback on
local code changes. (As a bonus, the mapping will also save the file for you.)

Of course, you can also change the key mapping to include whatever chef-runner
options you need. For example, I like using a Chef output format that is less
verbose:

```vim
nnoremap <leader>r :w\|!chef-runner -F min -l warn %<cr>
```

Last but not least, you can always reprogram the Vim key mapping at runtime if
you need to pass project-specific options like host or machine name:

```vim
:nnoremap <leader>r :w\|!chef-runner -F min -l warn -H example.local %<cr>
```

## More Tips

You can further speed up working with chef-runner by doing the following:

* Split up Chef recipes into smaller logical chunks and include those chunks
  using the `include_recipe` method ([good example][pr-recipes]).
* Enable [SSH connection sharing and persistence][ssh-speedup] to speed up
  repeated SSH connections.

## Testing

[![Build Status](https://travis-ci.org/mlafeldt/chef-runner.png?branch=master)](https://travis-ci.org/mlafeldt/chef-runner)

chef-runner comes with a couple of Go unit tests and Cucumber features that help
to ensure the tool works as expected. You can run all tests this way:

    $ ./script/test

## License and Author

Author:: Mathias Lafeldt (<mathias.lafeldt@gmail.com>)

Copyright:: 2013-2014, Mathias Lafeldt

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.

## Contributing

Please see `CONTRIBUTING.md` for details.


[Berkshelf]: http://berkshelf.com/
[blog post]: http://mlafeldt.github.io/blog/2014/01/telling-people-about-chef-runner/
[demo video]: http://vimeo.com/78769511
[Go]: http://golang.org/doc/install
[leader key]: http://usevim.com/2012/07/20/vim101-leader/
[pr-cookbook]: https://github.com/elm-city-craftworks/practicing-ruby-cookbook#readme
[pr-recipes]: https://github.com/elm-city-craftworks/practicing-ruby-cookbook/tree/master/recipes
[run list]: http://docs.opscode.com/essentials_node_object_run_lists.html
[ssh-speedup]: http://interrobeng.com/2013/08/25/speed-up-git-5x-to-50x/
[Vagrant]: http://vagrantup.com/
[VirtualBox]: https://www.virtualbox.org/

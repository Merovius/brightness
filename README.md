# brightness - control backlight brightness via sysfs

This is yet another tiny tool to control backlight brightness, after xbacklight
stopped working for me. It works on my machine, but I didn't put a lot of
effort into making it portable.

You can install it via

```
go get github.com/Merovius/brightness
# To allow non-root users to modify brightness
sudo chown root:root $GOPATH/bin/brightness; sudo chmod ug+s $GOPATH/bin/brightness
```

Usage:

```
brightness [-display <disp>] [<expr>]
```

Where `<expr>` is a string matching `[+-][0-9]+%?`.

* If a sign is given, the brightness level is increased/decreased accordingly.
  Otherwise it's set to the given value.
* If a percent sign is given, the level is interpreted relative to the maximum.
* If no expression is given, the current level is printed.

# License

Copyright 2018 Axel Wagner

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

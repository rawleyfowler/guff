# guff
The ultra lightweight HTTP client with 0 dependencies. Guff is just a little less than 200 lines,
and performs the majority of the work that a tool like cURL does, but without the "bloat" of supporting a thousand protocols.

Guff does one thing, and that's HTTP, and it does it damn well.

## How to install
```
git clone https://github.com/rawleyfowler/guff.git
cd guff
doas make install
```

With Guff installed, give it a test with the following:
```
guff www.gnu.org
```
If you see a HTML response dump, it works! If not, make an issue (and double check your internet connection).

What if i'm on Windows? Free yourself and use Linux or BSD!
## License
Guff is provided under the GNU General Public License version 3.0, please read the `LICENSE` file at the root of the project
for more information about your rights pertaining to this software. Guff is FREE SOFTWARE.
/*
 * You may redistribute this program and/or modify it under the terms of
 * the GNU General Public License as published by the Free Software Foundation,
 * either version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

"use strict";

function timeNow() {
  let time = new Date().toISOString();

  return time;
}

var log = {
  error: function(msg) {
    console.error('%s - \x1B[31mError\x1B[00m: %s', timeNow(), msg);
  },

  warn: function(msg) {
    console.warn('%s - \x1B[33mWarning:\x1B[00m %s', timeNow(), msg);
  },

  info: function(msg) {
    console.info('%s - Info: %s', timeNow(), msg);
  },

  debug: function(msg) {
    if (process.env['NODE_DEBUG'] && (process.env['NODE_DEBUG'].split(',')).indexOf('iptd') > -1) {
      console.log('%s - \x1B[34mDebug:\x1B[00m %s', timeNow(), msg);
    }
  }
};

module.exports = log;

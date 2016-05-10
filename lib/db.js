"use strict";

var fs   = require('fs'),
    path = require('path');

const DB_DIR_MODE = 511;

/**
 * Filesystem based database, it's a wrapper for fs that changes
 * root '/' to the database path.
 * @param dbPath  Database path
 */
var DB = function DB (dbPath) {
  this.path = path.normalize(dbPath);

  if(!fs.existsSync(dbPath)) {
    fs.mkdirSync(dbPath, DB_DIR_MODE);
  }
};

/**
 * Parse JSON and return in callback
 * @param   cb  callback(error, data)
 * @returns function(err, data)
 */
function parseJSON(cb) {
  return function (err, buffer) {
    if (err) {
      return cb(err);
    }

    try {
      return cb(null, JSON.parse(buffer.toString()));
    } catch (err) {
      return cb(err);
    }
  };
}

/**
 * Read from filesystem
 * @param   file      File to read from, will read every file if null
 * @param   callback  callback(error, data)
 * @returns function(error, data)
 */
DB.prototype.read = function read(file, callback) {
  let db = this;    

  if (file === null) {
    fs.readdir(db.path, function(err, files) {
      if (err) {
        return callback(err);
      }

      for(let file = 0; file < files.length; file++) {
        fs.readFile(db.path + '/' + files[file], parseJSON(callback));

      }
    });

  } else {

    fs.readFile(db.path + '/' + file, parseJSON(callback));
  }
};

/**
 * Check if file exists
 * @param   file      File to check if exists
 * @param   callback  callback(result)
 * @returns function(result)
 */
DB.prototype.exists = function exists(file, callback) {
  fs.exists(this.path + '/' + file, function(result) {
    return callback(result);
  });
};

/**
 * Write to filesystem
 * @param   file      File to write to
 * @param   callback  callback(error)
 * @returns function(error)
 */
DB.prototype.write = function write(file, data, callback) {
  data = JSON.stringify(data);

  fs.writeFile(this.path + '/' + file, data, function(err) {
    if (err) {
      return callback(err);
    }

    return callback();
  });
};

module.exports = DB;

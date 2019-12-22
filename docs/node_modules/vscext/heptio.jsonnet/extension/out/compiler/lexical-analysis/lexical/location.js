"use strict";
//////////////////////////////////////////////////////////////////////////////
// Location
Object.defineProperty(exports, "__esModule", { value: true });
// Location represents a single location in an (unspecified) file.
class Location {
    constructor(line, column) {
        this.line = line;
        this.column = column;
        // IsSet returns if this Location has been set.
        this.IsSet = () => {
            return this.line != 0;
        };
        this.String = () => {
            return `${this.line}:${this.column}`;
        };
        this.beforeRangeOrEqual = (range) => {
            const begin = range.begin;
            if (this.line < begin.line) {
                return true;
            }
            else if (this.line == begin.line && this.column <= begin.column) {
                return true;
            }
            return false;
        };
        this.strictlyBeforeRange = (range) => {
            const begin = range.begin;
            if (this.line < begin.line) {
                return true;
            }
            else if (this.line == begin.line && this.column < begin.column) {
                return true;
            }
            return false;
        };
        this.afterRangeOrEqual = (range) => {
            const end = range.end;
            if (this.line > end.line) {
                return true;
            }
            else if (this.line == end.line && this.column >= end.column) {
                return true;
            }
            return false;
        };
        this.strictlyAfterRange = (range) => {
            const end = range.end;
            if (this.line > end.line) {
                return true;
            }
            else if (this.line == end.line && this.column > end.column) {
                return true;
            }
            return false;
        };
        this.inRange = (loc) => {
            const range = {
                beginLine: loc.begin.line,
                endLine: loc.end.line,
                beginCol: loc.begin.column,
                endCol: loc.end.column,
            };
            if (range.beginLine == this.line && this.line == range.endLine &&
                range.beginCol <= this.column && this.column <= range.endCol) {
                return true;
            }
            else if (range.beginLine < this.line && this.line == range.endLine &&
                this.column <= range.endCol) {
                return true;
            }
            else if (range.beginLine == this.line && this.line < range.endLine &&
                this.column >= range.beginCol) {
                return true;
            }
            else if (range.beginLine < this.line && this.line < range.endLine) {
                return true;
            }
            else {
                return false;
            }
        };
    }
}
Location.fromString = (coord) => {
    const nums = coord.split(":");
    if (nums.length != 2) {
        return null;
    }
    return new Location(parseInt(nums[0]), parseInt(nums[1]));
};
exports.Location = Location;
const emptyLocation = () => new Location(0, 0);
//////////////////////////////////////////////////////////////////////////////
// LocationRange
// LocationRange represents a range of a source file.
class LocationRange {
    constructor(fileName, begin, end) {
        this.fileName = fileName;
        this.begin = begin;
        this.end = end;
        // IsSet returns if this LocationRange has been set.
        this.IsSet = () => {
            return this.begin.IsSet();
        };
        this.String = () => {
            if (!this.IsSet()) {
                return this.fileName;
            }
            let filePrefix = "";
            if (this.fileName.length > 0) {
                filePrefix = this.fileName + ":";
            }
            if (this.begin.line == this.end.line) {
                if (this.begin.column == this.end.column) {
                    return `${filePrefix}${this.begin.String()}`;
                }
                return `${filePrefix}${this.begin.String()}-${this.end.column}`;
            }
            return `${filePrefix}(${this.begin.String()})-(${this.end.String()})`;
        };
        this.rangeIsTighter = (thatRange) => {
            return this.begin.inRange(thatRange) && this.end.inRange(thatRange);
        };
    }
}
LocationRange.fromString = (filename, loc) => {
    // NOTE: Use `g` to search the string for all coordinates
    // formatted as `x:y`.
    const coordinates = loc.match(/(\d+:\d+)+/g);
    let start = null;
    let end = null;
    if (coordinates == null) {
        console.log(`Could not parse coordinates '${loc}'`);
        return null;
    }
    else if (coordinates.length == 2) {
        // Easy case. Of the form `(x1:y1)-(x2:y2)`.
        start = Location.fromString(coordinates[0]);
        end = Location.fromString(coordinates[1]);
        return start != null && end != null && new LocationRange(filename, start, end) || null;
    }
    else if (coordinates.length == 1) {
        // One of two forms: `x1:y1` or `x1:y1-y2`.
        start = Location.fromString(coordinates[0]);
        if (start == null) {
            return null;
        }
        const y2 = loc.match(/\-(\d+)/);
        if (y2 == null) {
            end = start;
        }
        else {
            end = new Location(start.line, parseInt(y2[1]));
        }
        return new LocationRange(filename, start, end);
    }
    else {
        console.log(`Could not parse coordinates '${loc}'`);
        return null;
    }
};
exports.LocationRange = LocationRange;
// This is useful for special locations, e.g. manifestation entry point.
exports.MakeLocationRangeMessage = (msg) => {
    return new LocationRange(msg, emptyLocation(), emptyLocation());
};
exports.MakeLocationRange = (fn, begin, end) => {
    return new LocationRange(fn, begin, end);
};
//# sourceMappingURL=location.js.map
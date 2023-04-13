declare global {
  interface Array<T> {
    /** @deprecated this only exists to help migrate Vector to Array */
    get(idx: number): T;
    /** @deprecated this only exists to help migrate Vector to Array */
    set(idx: number, value: T): void;
    /** @deprecated this only exists to help migrate Vector to Array */
    add(value: T): void;
    /** @deprecated this only exists to help migrate Vector to Array */
    toArray(): T[];
  }
}

// JS original sin
// writable: true because Jest will re-exec this block multiple times (in a browser this only runs once)
Object.defineProperties(Array.prototype, {
  get: {
    value: function(idx: number): any {
      return (this as any)[idx];
    },
    writable: true,
    enumerable: false,
    configurable: false,
  },
  set: {
    value: function(idx: number, value: any) {
      (this as any)[idx] = value;
    },
    writable: true,
    enumerable: false,
    configurable: false,
  },
  add: {
    value: function(value: any) {
      (this as any).push(value);
    },
    writable: true,
    enumerable: false,
    configurable: false,
  },
  toArray: {
    value: function() {
      return this;
    },
    writable: true,
    enumerable: false,
    configurable: false,
  },
});

/** @deprecated use a simple Array<T> */
export interface Vector<T = any> extends Array<T> {
  length: number;

  /**
   * Access the value by index (Like an array)
   */
  get(index: number): T;

  /**
   * Set a value
   */
  set: (index: number, value: T) => void;

  /**
   * Adds the value to the vector
   * Same as Array.push()
   */
  add: (value: T) => void;

  /**
   * Get the results as an array.
   */
  toArray(): T[];
}

/**
 * Apache arrow vectors are Read/Write
 *
 * @deprecated -- this is now part of the base Vector interface
 */
export interface ReadWriteVector<T = any> extends Vector<T> {}

/**
 * Vector with standard manipulation functions
 *
 * @deprecated -- this is now part of the base Vector interface
 */
export interface MutableVector<T = any> extends ReadWriteVector<T> {}

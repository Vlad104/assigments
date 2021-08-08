class Mutex {
  constructor() {
    this.queue = [];
    this.locked = false;
  }

  lock() {
    return new Promise((resolve) => {
        if (this.locked) {
            this.queue.push(resolve);
        } else {
            this.locked = true;
            resolve();
        }
    });
  }

  unlock() {
      if (this.queue.length > 0) {
          const resolve = this.queue.shift();
          resolve();
      } else {
          this.locked = false;
      }
  }
}

const mutex = new Mutex();

function runExclusively(func) {
  mutex.lock().then(() => {
    func().then(() => {
      mutex.unlock();
    })
})
}

async function someAsyncFunc() {
  await new Promise(resolve => {
    setTimeout(() => {
      resolve()
    } , 10)
  })
}

async function someExclusiveFunc(id) {
    console.log('start', id);
    await someAsyncFunc();
    console.log('end', id);
}

runExclusively(someExclusiveFunc.bind(this, 1));
runExclusively(someExclusiveFunc.bind(this, 2));
runExclusively(someExclusiveFunc.bind(this, 3));
runExclusively(someExclusiveFunc.bind(this, 4));

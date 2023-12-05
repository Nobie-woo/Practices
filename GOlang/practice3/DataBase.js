const express = require('express');
const app = express();
app.use(express.json());

class HashTable {
  constructor(size = 256) {
    this.size = size;
    this.table = new Array(size);
  }

  hash(key) {
    let hash = 0;
    for (let i = 0; i < key.length; i++) {
      hash = (hash * 31 + key.charCodeAt(i)) % this.size;
    }
    return hash;
  }

  insert(key, value) { // добавление в таблицу
    let index = this.hash(key);
    let originalIndex = index;
    while (this.table[index] !== undefined && this.table[index].key !== key) {
      index = (index + 1) % this.size;
      if (index === originalIndex) {
        throw new Error("Хеш-таблица переполнена");
      }
    }
    this.table[index] = { key, value };
  }

  getKeyByValue(value) { // проверка есть ли уже такая длинная ссылка
    for (let i = 0; i < this.size; i++) {
      if (this.table[i] !== undefined && this.table[i].value === value) {
        return this.table[i].key;
      }
    }
    return "";
  }

  get(key) { // получение длинной ссылки для перехода
    let index = this.hash(key);
    let originalIndex = index;
    while (this.table[index] !== undefined) {
      if (this.table[index].key === key) {
        return this.table[index].value; // нашли значение по ключу
      }
      index = (index + 1) % this.size; // переходим к следующей ячейке
      if (index === originalIndex) {
        return undefined; // возвращаем undefined, если ключ не найден после просмотра всех ячеек
      }
    }
    return undefined; // возвращаем undefined, если ключ не найден
  }
}

app.listen(8081, () => {
  console.log('Server listening on port 8081');
});

const hashTable = new HashTable();
  
  app.post('/insert', (req, res) => { // принятие пост-запроса от сервера на сокращение ссылки
    const key = req.body.key;
    const value = req.body.value;
    //console.log(key, " ", value, " ", hashTable.getKeyByValue(value));
    if (hashTable.getKeyByValue(value) == "") {
      hashTable.insert(key, value);
      res.send(key);
    }
    else {
      res.send(hashTable.getKeyByValue(value));
    }
  });
  
  app.get('/get', async (req, res) => { // получение длинной ссылки для отправки на сервер
    const key = req.query.key;
    //console.log(key);
    try {
      const response = await hashTable.get('http://localhost:8080/' + key);
      //console.log(key, " ", response);
      res.send(response);
    } catch (error) {
      res.status(500).send("Failed to retrieve the value");
    }
  });

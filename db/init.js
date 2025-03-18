db = db.getSiblingDB(env.MONGO_INITDB_DATABASE);
db.createCollection(env.MONGO_INITDB_COLLECTION);
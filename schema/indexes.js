db.campaigns.createIndex({ status: 1 });

db.campaigns.createIndex(
  { "targeting.includeCountry": 1 },
  { collation: { locale: "en", strength: 2 } }
);

db.campaigns.createIndex(
  { "targeting.includeOS": 1 },
  { collation: { locale: "en", strength: 2 } }
);

db.campaigns.createIndex(
  { "targeting.includeApp": 1 },
  { collation: { locale: "en", strength: 2 } }
);

db.campaigns.createIndex(
  { status: 1, "targeting.includeCountry": 1 },
  { collation: { locale: "en", strength: 2 } }
);
db.campaigns.createIndex(
  { status: 1, "targeting.includeApp": 1 },
  { collation: { locale: "en", strength: 2 } }
);
db.campaigns.createIndex(
  { status: 1, "targeting.includeOS": 1 },
  { collation: { locale: "en", strength: 2 } }
);

db.campaigns.find({
  status: "ACTIVE",
  $and: [
    {
      $or: [
        { "targeting.includeOS": "Android" },
        { "targeting.includeOS": { $size: 0 } },
      ],
    },
    {
      $or: [
        { "targeting.includeCountry": "Germany" },
        { "targeting.includeCountry": { $size: 0 } },
      ],
    },

    {
      $or: [
        { "targeting.includeApp": "com.jetrun.game" },
        { "targeting.includeApp": { $size: 0 } },
      ],
    },
  ],
  "targeting.excludeApp": { $nin: ["com.jetrun.game"] },
  "targeting.excludeCountry": { $nin: ["Germany"] },
  "targeting.excludeOS": { $nin: ["Android"] },
});

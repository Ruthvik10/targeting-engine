db.createCollection("campaigns", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["name", "image", "cta", "status", "targeting"],
      properties: {
        name: {
          bsonType: "string",
          description: "Campaign name must be a string and is required",
        },
        image: {
          bsonType: "string",
          description: "Image URL must be a string and is required",
        },
        cta: {
          bsonType: "string",
          description: "CTA (Call-to-Action) must be a string and is required",
        },
        status: {
          enum: ["ACTIVE", "INACTIVE"],
          description:
            "Status can only be either ACTIVE or INACTIVE and is required",
        },
        targeting: {
          bsonType: "object",
          required: [
            "includeApp",
            "excludeApp",
            "includeOS",
            "excludeOS",
            "includeCountry",
            "excludeCountry",
          ],
          properties: {
            includeApp: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of included App IDs",
            },
            excludeApp: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of excluded App IDs",
            },
            includeOS: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of included Operating Systems",
            },
            excludeOS: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of excluded Operating Systems",
            },
            includeCountry: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of included Countries",
            },
            excludeCountry: {
              bsonType: "array",
              items: {
                bsonType: "string",
              },
              description: "Array of excluded Countries",
            },
          },
          description: "Targeting rules for the campaign",
        },
      },
    },
  },
});

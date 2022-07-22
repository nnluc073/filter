### Settings

You can disable certain features, or blacklist certain fields using `filter.Settings`:

```go
settings := &filter.Settings{
	DisableFields: true, // Prevent usage of "fields"
	DisableFilter: true, // Prevent usage of "filter"
	DisableSort:   true, // Prevent usage of "sort"
	DisableJoin:   true, // Prevent usage of "join"

	FieldsSearch:   []string{"a", "b"},      // Optional, the fields used for the search feature
	SearchOperator: filter.Operators["$eq"], // Optional, operator used for the search feature, defaults to "$cont"

	Blacklist: filter.Blacklist{
		// Prevent selecting, sorting and filtering on these fields
		FieldsBlacklist: []string{"a", "b"},

		// Prevent joining these relations
		RelationsBlacklist: []string{"Relation"},

		Relations: map[string]*filter.Blacklist{
			// Blacklist settings to apply to this relation
			"Relation": &filter.Blacklist{
				FieldsBlacklist:    []string{"c", "d"},
				RelationsBlacklist: []string{"Parent"},
				Relations:          map[string]*filter.Blacklist{ /*...*/ },
				IsFinal:            true, // Prevent joining any child relation if true
			},
		},
	},
}
paginator, tx := settings.Scope(database.GetConnection(), request, &results)
```

### Filter

> ?filter=**field**||**$operator**||**value**

*Examples:*

> ?filter=**name**||**$cont**||**Jack** (`WHERE name LIKE "%Jack%"`)

You can add multiple filters. In that case, it is interpreted as an `AND` condition.

You can use `OR` conditions using `?or` instead, or in combination:

> ?filter=**name**||**$cont**||**Jack**&or=**name**||**$cont**||**John**  (`WHERE (name LIKE %Jack% OR name LIKE "%John%")`)  
> ?filter=**age**||**$eq**||**50**&filter=**name**||**$cont**||**Jack**&or=**name**||**$cont**||**John** (`WHERE ((age = 50 AND name LIKE "%Jack%") OR name LIKE "%John%")`)

You can filter using columns from one-to-one relations ("has one" or "belongs to"):

> ?filter=**Relation.name**||**$cont**||**Jack**

If there is only one "or", it is considered as a regular filter:

> ?or=**name**||**$cont**||**John**  (`WHERE name LIKE "%John%"`)  

If both "filter" and "or" are present, then they are interpreted as a combination of two `AND` groups compared with each other using `OR`:

> ?filter=**age**||**$eq**||**50**&filter=**name**||**$cont**||**Jack**&or=**name**||**$cont**||**John**&or=**name**||**$cont**||**Doe**  
> `WHERE ((age = 50 AND name LIKE "%Jack%") OR (name LIKE "%John%" AND name LIKE "%Doe%"))`

**Note:** All the filter conditions added to the SQL query are **grouped** (surrounded by parenthesis). 

#### Operators

|                |                                                         |
|----------------|---------------------------------------------------------|
| **`$eq`**      | `=`, equals                                             |
| **`$ne`**      | `<>`, not equals                                        |
| **`$gt`**      | `>`, greater than                                       |
| **`$lt`**      | `<`, lower than                                         |
| **`$gte`**     | `>=`, greater than or equals                            |
| **`$lte`**     | `<=`, lower than or equals                              |
| **`$starts`**  | `LIKE val%`, starts with                                |
| **`$ends`**    | `LIKE %val`, ends with                                  |
| **`$cont`**    | `LIKE %val%`, contains                                  |
| **`$excl`**    | `NOT LIKE %val%`, not contains                          |
| **`$in`**      | `IN (val1, val2,...)`, in (accepts multiple values)     |
| **`$notin`**   | `NOT IN (val1, val2,...)`, in (accepts multiple values) |
| **`$isnull`**  | `IS NULL`, is NULL (doesn't accept value)               |
| **`$notnull`** | `IS NOT NULL`, not NULL (doesn't accept value)          |
| **`$between`** | `BETWEEN val1 AND val2`, between (accepts two values)   |

### Search

Search is similar to multiple `or=column||$cont||value`, but the column and operator are specified by the server instead of the client.

Specify the column using `Settings`:
```go
settings := &filter.Settings{
	FieldsSearch: []string{"a", "b"},
	SearchOperator: filter.Operators["$eq"], // Optional, defaults to "$cont"
	//...
}
```

> ?search=John (`WHERE (a LIKE "%John%" OR b LIKE "%John%")`)

If you don't specify `FieldsSearch`, the query will search in all selectable fields.

### Fields / Select

> ?fields=**field1**,**field2**

A comma-separated list of fields to select. If this field isn't provided, uses `SELECT *`.

### Sort

> ?sort=**column**,**ASC**|**DESC**

*Examples:*

> ?sort=**name**,**ASC**  
> ?sort=**age**,**DESC**

You can also sort by multiple fields.

> ?sort=**age**,**DESC**&sort=**name**,**ASC**

### Join

> ?join=**relation**

Preload a relation. You can also only select the columns you need:

> ?join=**relation**||**field1**,**field2**

You can join multiple relations:

> ?join=**profile**||**firstName**,**email**&join=**notifications**||**content**&join=**tasks**

### Pagination

Internally, `goyave.dev/filter` uses [Goyave's `Paginator`](https://goyave.dev/guide/basics/database.html#pagination).

> ?page=**1**&per_page=**10**

- If `page` isn't given, the first page will be returned.
- If `per_page` isn't given, the default page size will be used. This default value can be overridden by changing `filter.DefaultPageSize`.
- Either way, the result is **always** paginated, even if those two parameters are missing.

## Security

- Inputs are escaped to prevent SQL injections.
- Fields are pre-processed and clients cannot request fields that don't exist. This prevents database errors. If a non-existing field is required, it is simply ignored. The same goes for sorts and joins. It is not possible to request a relation that doesn't exist.
- Foreign keys are always selected in joins to ensure associations can be assigned to parent model.
- **Be careful** with bidirectional relations (for example an article is written by a user, and a user can have many articles). If you enabled both your models to preload these relations, the client can request them with an infinite depth (`Articles.User.Articles.User...`). To prevent this, it is advised to use **the relation blacklist** or **IsFinal** on the deepest requestable models. See the settings section for more details.
# fiberfilter

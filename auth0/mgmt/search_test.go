package mgmt_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zenoss/go-auth0/auth0/http"
	mocks "github.com/zenoss/go-auth0/auth0/http/mocks"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)

const testJson = `{
  "start": 26,
  "limit": 26,
  "length": 26,
  "users": [
    {
      "created_at": "2020-03-30T14:09:02.845Z",
      "email": "joedoe5@gmail.com",
      "email_verified": false,
      "family_name": "doe5",
      "given_name": "joe5",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vqvlh7q0ts9vrp6o0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe5@gmail.com",
      "nickname": "joedoe5",
      "picture": "https://s.gravatar.com/avatar/c24fdbfc376c72ad673cf9bb5a5699ca?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:09:02.845Z",
      "user_id": "auth0|tenant:lysa2603:bq0vqvlh7q0ts9vrp6o0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:09:24.039Z",
      "email": "joedoe6@gmail.com",
      "email_verified": false,
      "family_name": "doe6",
      "given_name": "joe6",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vr4th7q0ts9vrp6og",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe6@gmail.com",
      "nickname": "joedoe6",
      "picture": "https://s.gravatar.com/avatar/85b814f484f6641c608896d72895247f?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:09:24.039Z",
      "user_id": "auth0|tenant:lysa2603:bq0vr4th7q0ts9vrp6og",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:09:45.351Z",
      "email": "joedoe7@gmail.com",
      "email_verified": false,
      "family_name": "doe7",
      "given_name": "joe7",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vradh7q0ts9vrp6p0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe7@gmail.com",
      "nickname": "joedoe7",
      "picture": "https://s.gravatar.com/avatar/43b75bc0d0be1d86907246dc0acd3e44?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:09:45.351Z",
      "user_id": "auth0|tenant:lysa2603:bq0vradh7q0ts9vrp6p0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:10:08.123Z",
      "email": "joedoe8@gmail.com",
      "email_verified": false,
      "family_name": "doe8",
      "given_name": "joe8",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vrfth7q0ts9vrp6pg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe8@gmail.com",
      "nickname": "joedoe8",
      "picture": "https://s.gravatar.com/avatar/d55a45491d6ee488190b96689e563b1e?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:10:08.123Z",
      "user_id": "auth0|tenant:lysa2603:bq0vrfth7q0ts9vrp6pg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:10:39.237Z",
      "email": "joedoe9@gmail.com",
      "email_verified": false,
      "family_name": "doe9",
      "given_name": "joe9",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vrnth7q0ts9vrp6q0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe9@gmail.com",
      "nickname": "joedoe9",
      "picture": "https://s.gravatar.com/avatar/b88845ab250b5f27fb8bafa6d8bfd84f?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:10:39.237Z",
      "user_id": "auth0|tenant:lysa2603:bq0vrnth7q0ts9vrp6q0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:11:00.988Z",
      "email": "joedoe10@gmail.com",
      "email_verified": false,
      "family_name": "doe10",
      "given_name": "joe10",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vrt5h7q0ts9vrp6qg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "joedoe10@gmail.com",
      "nickname": "joedoe10",
      "picture": "https://s.gravatar.com/avatar/c539c4ddc6cb156ba7273f176f5f63bd?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
      "updated_at": "2020-03-30T14:11:00.988Z",
      "user_id": "auth0|tenant:lysa2603:bq0vrt5h7q0ts9vrp6qg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:11:47.930Z",
      "email": "samdeen1@gmail.com",
      "email_verified": false,
      "family_name": "deen1",
      "given_name": "sam1",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vs8th7q0ts9vrp6r0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen1@gmail.com",
      "nickname": "samdeen1",
      "picture": "https://s.gravatar.com/avatar/1bc82ac1008f58322502db1838fd335b?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:11:47.930Z",
      "user_id": "auth0|tenant:lysa2603:bq0vs8th7q0ts9vrp6r0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:12:14.560Z",
      "email": "samdeen2@gmail.com",
      "email_verified": false,
      "family_name": "deen2",
      "given_name": "sam2",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vsflh7q0ts9vrp6rg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen2@gmail.com",
      "nickname": "samdeen2",
      "picture": "https://s.gravatar.com/avatar/b19c6a3e01279a825bcd9f36a8deb455?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:12:14.560Z",
      "user_id": "auth0|tenant:lysa2603:bq0vsflh7q0ts9vrp6rg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:12:40.320Z",
      "email": "samdeen3@gmail.com",
      "email_verified": false,
      "family_name": "deen3",
      "given_name": "sam3",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vsm5h7q0ts9vrp6s0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen3@gmail.com",
      "nickname": "samdeen3",
      "picture": "https://s.gravatar.com/avatar/16b7b0a63f2db5b832e9b664165a65d3?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:12:40.320Z",
      "user_id": "auth0|tenant:lysa2603:bq0vsm5h7q0ts9vrp6s0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:13:09.431Z",
      "email": "samdeen4@gmail.com",
      "email_verified": false,
      "family_name": "deen4",
      "given_name": "sam4",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vstdh7q0ts9vrp6sg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen4@gmail.com",
      "nickname": "samdeen4",
      "picture": "https://s.gravatar.com/avatar/a85a70309acff14787b23c538acf8f66?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:13:09.431Z",
      "user_id": "auth0|tenant:lysa2603:bq0vstdh7q0ts9vrp6sg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:13:40.071Z",
      "email": "samdeen5@gmail.com",
      "email_verified": false,
      "family_name": "deen5",
      "given_name": "sam5",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vt4th7q0ts9vrp6t0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen5@gmail.com",
      "nickname": "samdeen5",
      "picture": "https://s.gravatar.com/avatar/05bb86f81d092b12823aa2711334adb0?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:13:40.071Z",
      "user_id": "auth0|tenant:lysa2603:bq0vt4th7q0ts9vrp6t0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:14:05.625Z",
      "email": "samdeen6@gmail.com",
      "email_verified": false,
      "family_name": "deen6",
      "given_name": "sam6",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vtbdh7q0ts9vrp6tg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen6@gmail.com",
      "nickname": "samdeen6",
      "picture": "https://s.gravatar.com/avatar/019b8821200afd096a101b912c00f16a?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:14:05.625Z",
      "user_id": "auth0|tenant:lysa2603:bq0vtbdh7q0ts9vrp6tg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:14:35.637Z",
      "email": "samdeen7@gmail.com",
      "email_verified": false,
      "family_name": "denn7",
      "given_name": "sam7",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vtith7q0ts9vrp6u0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen7@gmail.com",
      "nickname": "samdeen7",
      "picture": "https://s.gravatar.com/avatar/535b530325b46c34f4c95fb28cf5b097?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:14:35.637Z",
      "user_id": "auth0|tenant:lysa2603:bq0vtith7q0ts9vrp6u0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:15:03.104Z",
      "email": "samdeen8@gmail.com",
      "email_verified": false,
      "family_name": "deen8",
      "given_name": "sam8",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vtplh7q0ts9vrp6ug",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen8@gmail.com",
      "nickname": "samdeen8",
      "picture": "https://s.gravatar.com/avatar/724a2ffa9c7b08714a73d7dfc0b1f4cf?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:15:03.104Z",
      "user_id": "auth0|tenant:lysa2603:bq0vtplh7q0ts9vrp6ug",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:15:32.947Z",
      "email": "samdeen9@gmail.com",
      "email_verified": false,
      "family_name": "deen9",
      "given_name": "sam9",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vu15h7q0ts9vrp6v0",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen9@gmail.com",
      "nickname": "samdeen9",
      "picture": "https://s.gravatar.com/avatar/41e3ff75b3fe29c8af2b064754aca1d0?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:15:32.947Z",
      "user_id": "auth0|tenant:lysa2603:bq0vu15h7q0ts9vrp6v0",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:15:58.061Z",
      "email": "samdeen10@gmail.com",
      "email_verified": false,
      "family_name": "deen10",
      "given_name": "sam10",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vu7dh7q0ts9vrp6vg",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "samdeen10@gmail.com",
      "nickname": "samdeen10",
      "picture": "https://s.gravatar.com/avatar/f6f83c2b2cbefdeabe876911444e354b?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsa.png",
      "updated_at": "2020-03-30T14:15:58.061Z",
      "user_id": "auth0|tenant:lysa2603:bq0vu7dh7q0ts9vrp6vg",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:16:28.931Z",
      "email": "timbird1@gmail.com",
      "email_verified": false,
      "family_name": "bird1",
      "given_name": "tim1",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vuf5h7q0ts9vrp700",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird1@gmail.com",
      "nickname": "timbird1",
      "picture": "https://s.gravatar.com/avatar/a829a18b96c4b57ed53cbc07ba8e67a6?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:16:28.931Z",
      "user_id": "auth0|tenant:lysa2603:bq0vuf5h7q0ts9vrp700",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:16:51.695Z",
      "email": "timbird2@gmail.com",
      "email_verified": false,
      "family_name": "bird2",
      "given_name": "tim2",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vukth7q0ts9vrp70g",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird2@gmail.com",
      "nickname": "timbird2",
      "picture": "https://s.gravatar.com/avatar/24a624261eee3e8ea14873deaf5fc597?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:16:51.695Z",
      "user_id": "auth0|tenant:lysa2603:bq0vukth7q0ts9vrp70g",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:17:13.617Z",
      "email": "timbird3@gmail.com",
      "email_verified": false,
      "family_name": "bird3",
      "given_name": "tim3",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vuqdh7q0ts9vrp710",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird3@gmail.com",
      "nickname": "timbird3",
      "picture": "https://s.gravatar.com/avatar/12153c17459163db39fb222650a5b1a9?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:17:13.617Z",
      "user_id": "auth0|tenant:lysa2603:bq0vuqdh7q0ts9vrp710",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:17:50.863Z",
      "email": "timbird4@gmail.com",
      "email_verified": false,
      "family_name": "bird4",
      "given_name": "tim4",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vv3lh7q0ts9vrp71g",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird4@gmail.com",
      "nickname": "timbird4",
      "picture": "https://s.gravatar.com/avatar/376e4c1bd664ac6730b9db3e94996422?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:17:50.863Z",
      "user_id": "auth0|tenant:lysa2603:bq0vv3lh7q0ts9vrp71g",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:18:16.591Z",
      "email": "timbird5@gmail.com",
      "email_verified": false,
      "family_name": "bird5",
      "given_name": "tim5",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vva5h7q0ts9vrp720",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird5@gmail.com",
      "nickname": "timbird5",
      "picture": "https://s.gravatar.com/avatar/fc60e98420d27504bb56c293f83a3bd5?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:18:16.591Z",
      "user_id": "auth0|tenant:lysa2603:bq0vva5h7q0ts9vrp720",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:18:38.844Z",
      "email": "timbird6@gmail.com",
      "email_verified": false,
      "family_name": "bird6",
      "given_name": "tim6",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vvflh7q0ts9vrp72g",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird6@gmail.com",
      "nickname": "timbird6",
      "picture": "https://s.gravatar.com/avatar/1a18246e6bab3f26a9b1e697df9d101a?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:18:38.844Z",
      "user_id": "auth0|tenant:lysa2603:bq0vvflh7q0ts9vrp72g",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:19:06.443Z",
      "email": "timbird7@gmail.com",
      "email_verified": false,
      "family_name": "bird7",
      "given_name": "tim7",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vvmlh7q0ts9vrp730",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird7@gmail.com",
      "nickname": "timbird7",
      "picture": "https://s.gravatar.com/avatar/8a1b69375e40a7e4c9b00b13ec14b0ec?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:19:06.443Z",
      "user_id": "auth0|tenant:lysa2603:bq0vvmlh7q0ts9vrp730",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:19:25.338Z",
      "email": "timbird8@gmail.com",
      "email_verified": false,
      "family_name": "bird8",
      "given_name": "tim8",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq0vvrdh7q0ts9vrp73g",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird8@gmail.com",
      "nickname": "timbird8",
      "picture": "https://s.gravatar.com/avatar/7861c44d920b84762043f54569d9c2d9?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:19:25.338Z",
      "user_id": "auth0|tenant:lysa2603:bq0vvrdh7q0ts9vrp73g",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:19:46.772Z",
      "email": "timbird9@gmail.com",
      "email_verified": false,
      "family_name": "bird9",
      "given_name": "tim9",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq1000lh7q0ts9vrp740",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird9@gmail.com",
      "nickname": "timbird9",
      "picture": "https://s.gravatar.com/avatar/ff22ecfced7bba7b1da6a385a6150cc1?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:19:46.772Z",
      "user_id": "auth0|tenant:lysa2603:bq1000lh7q0ts9vrp740",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    },
    {
      "created_at": "2020-03-30T14:20:30.354Z",
      "email": "timbird10@gmail.com",
      "email_verified": false,
      "family_name": "bird10",
      "given_name": "tim10",
      "identities": [
        {
          "user_id": "tenant:lysa2603:bq100blh7q0ts9vrp74g",
          "connection": "Username-Password-Authentication",
          "provider": "auth0",
          "isSocial": false
        }
      ],
      "name": "timbird10@gmail.com",
      "nickname": "timbird10",
      "picture": "https://s.gravatar.com/avatar/ec5735e86f79586e9a4d87a97d84e326?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fti.png",
      "updated_at": "2020-03-30T14:20:30.354Z",
      "user_id": "auth0|tenant:lysa2603:bq100blh7q0ts9vrp74g",
      "app_metadata": {
        "tenant": "lysa2603"
      }
    }
  ],
  "total": 52
}`

const testJsonUsers = `
[
  {
    "created_at": "2020-03-30T14:09:02.845Z",
    "email": "joedoe5@gmail.com",
    "email_verified": false,
    "family_name": "doe5",
    "given_name": "joe5",
    "identities": [
      {
        "user_id": "tenant:lysa2603:bq0vqvlh7q0ts9vrp6o0",
        "connection": "Username-Password-Authentication",
        "provider": "auth0",
        "isSocial": false
      }
    ],
    "name": "joedoe5@gmail.com",
    "nickname": "joedoe5",
    "picture": "https://s.gravatar.com/avatar/c24fdbfc376c72ad673cf9bb5a5699ca?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fjo.png",
    "updated_at": "2020-03-30T14:09:02.845Z",
    "user_id": "auth0|tenant:lysa2603:bq0vqvlh7q0ts9vrp6o0",
    "app_metadata": {
      "tenant": "lysa2603"
    }
  }
]
`

func TestUserSearch(t *testing.T) {
	mockDoer := &mocks.Doer{}
	mockHTTPRequest := mock.AnythingOfType("*http.Request")
	svc := mgmt.New(&http.Client{
		Doer: mockDoer,
	})
	mockDoer.On("Do", mockHTTPRequest, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			respBody, ok := args.Get(1).(*mgmt.UsersPage)
			if !ok {
				t.Errorf("response body should be of type %t", reflect.TypeOf(&mgmt.UsersPage{}))
			}
			err := json.Unmarshal([]byte(testJson), respBody)
			assert.NoErrorf(t, err, "failed to unmarshal JSON %q", err)
		}).Once()
	// usersPage, err := svc.Users.Search(opts)
	opts := mgmt.SearchUsersOpts{
		Q:             fmt.Sprintf(`app_metadata.tenant:"%s" AND identities.connection:"%s"`, "lysa2603", "Username-Password-Authentication"),
		SearchEngine:  "v3",
		IncludeTotals: true,
		Page:          0,
		PerPage:       26,
	}
	resp, err := svc.Users.Search(opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 26, resp.Start)
	assert.Equal(t, 26, resp.Length)
	assert.Equal(t, 52, resp.Total)
	assert.NotEmpty(t, resp.Users)
	mockDoer.On("Do", mockHTTPRequest, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			respBody, ok := args.Get(1).(*[]mgmt.User)
			if !ok {
				t.Errorf("response body should be of type %t", reflect.TypeOf(&[]mgmt.User{}))
			}
			err := json.Unmarshal([]byte(testJsonUsers), respBody)
			assert.NoErrorf(t, err, "failed to unmarshal JSON %q", err)
		}).Once()

	opts.IncludeTotals = false
	resp2, err := svc.Users.Search(opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
	assert.NotEmpty(t, resp2.Users)
}

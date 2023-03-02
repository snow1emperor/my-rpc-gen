syntax = "proto3";

option go_package = "{{.packageName}}";

package {{.packageName}};

import "gogo.proto";


enum TLConstructor {
	UNKNOWN = 0;
{{.constructors}}
}

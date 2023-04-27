package utils

import (
    "bytes"
    "encoding/json"
)

func Prettyprint(jsond string) (string, error) {
    var pretty bytes.Buffer
    if err := json.Indent(&pretty, []byte(jsond), "", "   "); err != nil {
        return "", err
    }
    return pretty.String(), nil
}

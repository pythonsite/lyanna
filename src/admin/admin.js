import UIkit from "./base"

let $error = $('meta[name=errors]').attr('content');

if ($error) {
    UIkit.notification({
        message: $error,
        status: 'danger',
        timeout: 1000
    })
}

import 'select2/dist/js/select2';
import 'select2/dist/css/select2.css';
import './admin';
import UIkit from './base';
import "../scss/select.scss";

let $switchInput = $('.switch-input');
let $switcher = $('.uk-switch input');

$switcher.on('click', (event)  => {
    let $this = $(event.currentTarget)[0];
    let checked = $this.checked;
    if (checked) {
        $this.setAttribute('value', 'on');
        $switchInput.attr('value', 'on');
    } else {
        $this.setAttribute('value', 'off');
        $switchInput.attr('value', 'off');
    }
});

$(document).ready(() => {
    $("select").select2({
        tags: true
    });
});
import 'select2/dist/js/select2';
import 'select2/dist/css/select2.css';
import './admin';
import UIkit from './base';
import "../scss/select.scss";

import SimpleMDE from '../vendor/simplemde';
import '../css/font-awesome.min.css';
import '../css/simplemde.min.css';
import '../scss/coremirror.scss';




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

let simplemde = new SimpleMDE({
    element: $("#content")[0],
    autoDownloadFontAwesome:false,//true从默认地址引入fontawesome依赖 false需自行引入(国内用bootcdn更快点)
    autofocus:true,
    autosave: {
        enabled: true,
        uniqueId: "SimpleMDE",
        delay: 1000,
    },
    blockStyles: {
        bold: "**",
        italic: "*",
        code: "```"
    },
    forceSync: true,
    hideIcons: false,
    indentWithTabs: true,
    lineWrapping: true,
    renderingConfig:{
        singleLineBreaks: false,
        codeSyntaxHighlighting: true // 需要highlight依赖
    },
    showIcons: true,
    spellChecker: true
});

let $content = $('meta[name=raw_content]').attr('content');
if ($content) {
    simplemde.value($content);
}

$(document).ready(() => {
    $("select").select2({
        tags: true
    });
});
let $commentContainer = $('.gitment-comments-list');
let $editorTab = $('.gitment-editor-tab');
let $editorPreview = $('.gitment-editor-preview')
let $editorWriteField = $('.gitment-editor-write-field');
let $writeTextarea = $editorWriteField.find('textarea');
let $editorPreviewField = $('.gitment-editor-preview-field');
let $loginBtn = $('.gitment-editor-login-link');
let $submitBtn = $('.gitment-editor-submit');


const target_id = $('meta[name=post_id]').attr('content');



$editorTab.click((e)=> {
    let self = $(e.currentTarget);
    $editorTab.removeClass('gitment-selected');
    self.addClass('gitment-selected');
    if (self.hasClass('preview')) {
        $editorWriteField.addClass('gitment-hidden')
        $editorPreviewField.removeClass('gitment-hidden')
        let text = $writeTextarea.val();
        if (!$writeTextarea.disabled) {
            if (!text) {
                $editorPreview.html('空空如也')
                return
            }
            $editorPreview.html('渲染中...')
            $.ajax({
                url: '/j/markdown',
                type: 'post',
                data: {'text': text},
                dataType: 'json',
                success: function (rs) {
                    $editorPreview.html(rs.text);
                }
            });
        }
    } else {
        $editorWriteField.removeClass('gitment-hidden')
        $editorPreviewField.addClass('gitment-hidden')
    }
});

$submitBtn.click((e)=> {
    let content = $writeTextarea.val();
    if (!content) {
        return
    }
    let self = $(e.currentTarget);
    self.html('提交...')
    self.attr('disabled', true)
    $.ajax({
        url: `/comment/${target_id}`,
        type: 'post',
        data: {'content': content},
        dataType: 'json',
        success: function (rs) {
            if (!rs.r) {
                $writeTextarea.val('')
                self.removeAttr('disabled')
                self.html('评论')
                $commentContainer.prepend(rs.html)
                console.log('评论成功')
            } else {
                console.log('评论失败')
            }
        }
    });
});
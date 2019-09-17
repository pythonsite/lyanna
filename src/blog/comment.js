let $commentContainer = $('.gitment-comments-list');
let $editorTab = $('.gitment-editor-tab');
let $editorPreview = $('.gitment-editor-preview')
let $editorWriteField = $('.gitment-editor-write-field');
let $writeTextarea = $editorWriteField.find('textarea');
let $editorPreviewField = $('.gitment-editor-preview-field');
let $loginBtn = $('.gitment-editor-login-link');
let $submitBtn = $('.gitment-editor-submit');
let $isEmptyDiv = $('.gitment-comments-empty');
let $pageItemBtn = $('.gitment-comments-page-item')


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
                url: '/comment/markdown',
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
        url: `/comment/post/${target_id}`,
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
                $isEmptyDiv.remove()
            } else {
                console.log('评论失败')
            }
        }
    });
});

// 翻页
$pageItemBtn.click((e)=> {
    let self = $(e.currentTarget), page;
    if (self.hasClass('gitment-selected')) {
        return
    }
    let current_page = parseInt($('.gitment-comments-pagination .gitment-selected').html())
    if (self.hasClass('next')) {
        page = current_page + 1
    } else if (self.hasClass('prev')) {
        page = current_page - 1
    } else {
        page = parseInt(self.html())
    }
    if (page != 1) {
        $pageItemBtn.eq(0).removeClass('gitment-hidden')
    } else {
        $pageItemBtn.eq(0).addClass('gitment-hidden')
    }
    if ($pageItemBtn.length - 2 <= page) {
        $pageItemBtn.eq(-1).addClass('gitment-hidden')
    } else {
        $pageItemBtn.eq(-1).removeClass('gitment-hidden')
    }
    $pageItemBtn.removeClass('gitment-selected');
    $pageItemBtn.eq(page).addClass('gitment-selected')
    const loading = document.createElement('div')
    loading.innerText = '加载评论...'
    loading.className = 'gitment-comments-loading'
    $commentContainer.empty().append(loading)
    $.ajax({
        url: `/comments/post/${target_id}?page=${page}&per_page=10`,
        type: 'get',
        dataType: 'json',
        success: function (rs) {
            $commentContainer.empty().append(rs.html)
        }
    })
});
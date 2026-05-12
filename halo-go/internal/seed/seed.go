package seed

import (
	"context"
	"log"
	"time"

	"github.com/halo-dev/halo-go/internal/data"
	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

func SeedData(store *data.ExtensionStore, username string) error {
	ctx := context.Background()
	now := time.Now()
	categoryName := "76514a40-6ef1-4ed9-b58a-e26945bde3ca"
	tagName := "c33ceabb-d8f1-4711-8991-bb8f5c92ad7c"
	postName := "5152aea5-c2e8-4717-8bba-2263d46e19d5"
	postSnapshotName := "fb5cd6bd-998d-4ccc-984d-5cc23b0a09f9"
	singlePageName := "373a5f79-f44f-441a-9df1-85a4f553ece8"
	pageSnapshotName := "c3f73cc2-194e-4cd8-9092-7386aa50a0e5"
	menuItemHome := "88c3f10b-321c-4092-86a8-70db00251b74"
	menuItemPost := "c4c814d1-0c2c-456b-8c96-4864965fee94"
	menuItemTag := "35869bd3-33b5-448b-91ee-cf6517a59644"
	menuItemPage := "b0d041fa-dc99-48f6-a193-8604003379cf"
	category := &model.Category{}
	category.Metadata.Name = categoryName
	category.Spec.DisplayName = "默认分类"
	category.Spec.Slug = "default"
	category.Spec.Description = "这是你的默认分类，如不需要，删除即可。"
	category.Status.Permalink = "/categories/default"
	if _, err := store.Create(ctx, category); err != nil {
		return err
	}
	tag := &model.Tag{}
	tag.Metadata.Name = tagName
	tag.Spec.DisplayName = "Halo"
	tag.Spec.Slug = "halo"
	tag.Status.Permalink = "/tags/halo"
	if _, err := store.Create(ctx, tag); err != nil {
		return err
	}
	postSnapshot := &model.Snapshot{}
	postSnapshot.Metadata.Name = postSnapshotName
	postSnapshot.Metadata.Annotations = map[string]string{"content.halo.run/keep-raw": "true"}
	postSnapshot.Spec.SubjectRef = extension.Ref{Group: "content.halo.run", Version: "v1alpha1", Kind: "Post", Name: postName}
	postSnapshot.Spec.RawType = "HTML"
	postSnapshot.Spec.RawPatch = `<h2 id="hello-halo"><strong>Hello Halo</strong></h2><p>如果你看到了这一篇文章，那么证明你已经安装成功了，感谢使用 <a target="_blank" rel="noopener noreferrer nofollow" href="https://www.halo.run/">Halo</a> 进行创作，希望能够使用愉快。</p><h2 id="%E7%9B%B8%E5%85%B3%E9%93%BE%E6%8E%A5"><strong>相关链接</strong></h2><ul><li><p>官网：<a target="_blank" rel="noopener noreferrer nofollow" href="https://www.halo.run">https://www.halo.run</a></p></li><li><p>文档：<a target="_blank" rel="noopener noreferrer nofollow" href="https://docs.halo.run">https://docs.halo.run</a></p></li><li><p>社区：<a target="_blank" rel="noopener noreferrer nofollow" href="https://bbs.halo.run">https://bbs.halo.run</a></p></li><li><p>应用市场：<a target="_blank" rel="noopener noreferrer nofollow" href="https://www.halo.run/store/apps">https://www.halo.run/store/apps</a></p></li><li><p>开源地址：<a target="_blank" rel="noopener noreferrer nofollow" href="https://github.com/halo-dev/halo">https://github.com/halo-dev/halo</a></p></li></ul><p>在使用过程中，有任何问题都可以通过以上链接找寻答案，或者联系我们。</p><blockquote><p>这是一篇自动生成的文章，请删除这篇文章之后开始你的创作吧！</p></blockquote>`
	postSnapshot.Spec.ContentPatch = postSnapshot.Spec.RawPatch
	postSnapshot.Spec.LastModifyTime = now
	postSnapshot.Spec.Owner = username
	postSnapshot.Spec.Contributors = []string{username}
	if _, err := store.Create(ctx, postSnapshot); err != nil {
		return err
	}
	post := &model.Post{}
	post.Metadata.Name = postName
	post.Spec.Title = "Hello Halo"
	post.Spec.Slug = "hello-halo"
	post.Spec.ReleaseSnapshot = postSnapshotName
	post.Spec.HeadSnapshot = postSnapshotName
	post.Spec.BaseSnapshot = postSnapshotName
	post.Spec.Owner = username
	post.Spec.Publish = true
	post.Spec.PublishTime = &now
	post.Spec.AllowComment = true
	post.Spec.Visible = "PUBLIC"
	post.Spec.Categories = []string{categoryName}
	post.Spec.Tags = []string{tagName}
	post.Spec.Excerpt.AutoGenerate = false
	post.Spec.Excerpt.Raw = "如果你看到了这一篇文章，那么证明你已经安装成功了，感谢使用 Halo 进行创作，希望能够使用愉快。"
	post.Status.Permalink = "/archives/hello-halo"
	if _, err := store.Create(ctx, post); err != nil {
		return err
	}
	pageSnapshot := &model.Snapshot{}
	pageSnapshot.Metadata.Name = pageSnapshotName
	pageSnapshot.Metadata.Annotations = map[string]string{"content.halo.run/keep-raw": "true"}
	pageSnapshot.Spec.SubjectRef = extension.Ref{Group: "content.halo.run", Version: "v1alpha1", Kind: "SinglePage", Name: singlePageName}
	pageSnapshot.Spec.RawType = "HTML"
	pageSnapshot.Spec.RawPatch = `<h2><strong>关于页面</strong></h2><p>这是一个自定义页面，你可以在后台的 <code>页面</code> -&gt; <code>自定义页面</code> 找到它，你可以用于新建关于页面、联系我们页面等等。</p><blockquote><p>这是一篇自动生成的页面，你可以在后台删除它。</p></blockquote>`
	pageSnapshot.Spec.ContentPatch = pageSnapshot.Spec.RawPatch
	pageSnapshot.Spec.LastModifyTime = now
	pageSnapshot.Spec.Owner = username
	pageSnapshot.Spec.Contributors = []string{username}
	if _, err := store.Create(ctx, pageSnapshot); err != nil {
		return err
	}
	singlePage := &model.SinglePage{}
	singlePage.Metadata.Name = singlePageName
	singlePage.Spec.Title = "关于"
	singlePage.Spec.Slug = "about"
	singlePage.Spec.Owner = username
	singlePage.Spec.Publish = true
	singlePage.Spec.AllowComment = true
	singlePage.Spec.Visible = "PUBLIC"
	singlePage.Spec.BaseSnapshot = pageSnapshotName
	singlePage.Spec.HeadSnapshot = pageSnapshotName
	singlePage.Spec.ReleaseSnapshot = pageSnapshotName
	singlePage.Spec.Excerpt.AutoGenerate = false
	singlePage.Spec.Excerpt.Raw = "这是一个自定义页面，你可以在后台 页面 -> 自定义页面 找到它，你可以用于新建关于页面、联系我们页面等等。"
	singlePage.Status.Permalink = "/about"
	if _, err := store.Create(ctx, singlePage); err != nil {
		return err
	}
	menuItemHomeObj := &model.MenuItem{}
	menuItemHomeObj.Metadata.Name = menuItemHome
	menuItemHomeObj.Spec.DisplayName = "首页"
	menuItemHomeObj.Spec.Href = "/"
	menuItemHomeObj.Spec.Priority = 0
	if _, err := store.Create(ctx, menuItemHomeObj); err != nil {
		return err
	}
	menuItemPostObj := &model.MenuItem{}
	menuItemPostObj.Metadata.Name = menuItemPost
	menuItemPostObj.Spec.DisplayName = "Hello Halo"
	menuItemPostObj.Spec.Href = "/archives/hello-halo"
	menuItemPostObj.Spec.Priority = 1
	menuItemPostObj.Spec.TargetRef = &extension.Ref{Group: "content.halo.run", Version: "v1alpha1", Kind: "Post", Name: postName}
	if _, err := store.Create(ctx, menuItemPostObj); err != nil {
		return err
	}
	menuItemTagObj := &model.MenuItem{}
	menuItemTagObj.Metadata.Name = menuItemTag
	menuItemTagObj.Spec.DisplayName = "Halo"
	menuItemTagObj.Spec.Href = "/tags/halo"
	menuItemTagObj.Spec.Priority = 2
	menuItemTagObj.Spec.TargetRef = &extension.Ref{Group: "content.halo.run", Version: "v1alpha1", Kind: "Tag", Name: tagName}
	if _, err := store.Create(ctx, menuItemTagObj); err != nil {
		return err
	}
	menuItemPageObj := &model.MenuItem{}
	menuItemPageObj.Metadata.Name = menuItemPage
	menuItemPageObj.Spec.DisplayName = "关于"
	menuItemPageObj.Spec.Href = "/about"
	menuItemPageObj.Spec.Priority = 3
	menuItemPageObj.Spec.TargetRef = &extension.Ref{Group: "content.halo.run", Version: "v1alpha1", Kind: "SinglePage", Name: singlePageName}
	if _, err := store.Create(ctx, menuItemPageObj); err != nil {
		return err
	}
	menu := &model.Menu{}
	menu.Metadata.Name = "primary"
	menu.Spec.DisplayName = "主菜单"
	menu.Spec.MenuItems = []string{menuItemHome, menuItemPost, menuItemTag, menuItemPage}
	if _, err := store.Create(ctx, menu); err != nil {
		return err
	}
	log.Println("初始化数据完成：默认分类、默认标签、Hello Halo 文章、关于页面、主菜单（含4个菜单项）")
	return nil
}

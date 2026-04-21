package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hashiURL = "https://releases.hashicorp.com/terraform/"

	hashicorpBody = `
	<li>
	<a href="/terraform/0.12.3-beta1/">terraform_0.12.3-beta1</a>
	</li>
	<li>
	<a href="/terraform/0.12.2/">terraform_0.12.2</a>
	</li>
	<li>
	<a href="/terraform/0.12.1/">terraform_0.12.1</a>
	</li>
	<li>
	<a href="/terraform/0.12.0/">terraform_0.12.0</a>
	</li>
	<li>
	<a href="/terraform/0.12.0-rc1/">terraform_0.12.0-rc1</a>
	</li>
	<li>
	<a href="/terraform/0.12.0-beta2/">terraform_0.12.0-beta2</a>
	</li>
	<li>
	<a href="/terraform/0.11.13/">terraform_0.11.13</a>
	</li>
`

	openTofuBody = `
<!DOCTYPE html>
<html>
<head>
    <title>OpenTofu releases</title>
</head>
<body>
<ul><li><a href="/tofu/1.7.1-beta1/">tofu_1.7.1-beta1</a></li><li><a href="/tofu/1.7.0/">tofu_1.7.0</a></li><li><a href="/tofu/1.7.0-rc1/">tofu_1.7.0-rc1</a></li><li><a href="/tofu/1.7.0-beta1/">tofu_1.7.0-beta1</a></li><li><a href="/tofu/1.7.0-alpha1/">tofu_1.7.0-alpha1</a></li><li><a href="/tofu/1.6.2/">tofu_1.6.2</a></li><li><a href="/tofu/1.6.0-alpha1/">tofu_1.6.0-alpha1</a></li></ul>
</body>
</html>
`

	hashicorpJSONData = `{"name":"terraform","versions":{"0.11.13":{"builds":[{"arch":"amd64","filename":"terraform_0.11.13_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_darwin_amd64.zip","version":"0.11.13"},{"arch":"arm64","filename":"terraform_0.11.13_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_darwin_arm64.zip","version":"0.11.13"},{"arch":"386","filename":"terraform_0.11.13_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_freebsd_386.zip","version":"0.11.13"},{"arch":"amd64","filename":"terraform_0.11.13_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_freebsd_amd64.zip","version":"0.11.13"},{"arch":"arm","filename":"terraform_0.11.13_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_freebsd_arm.zip","version":"0.11.13"},{"arch":"386","filename":"terraform_0.11.13_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_386.zip","version":"0.11.13"},{"arch":"amd64","filename":"terraform_0.11.13_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip","version":"0.11.13"},{"arch":"arm","filename":"terraform_0.11.13_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_arm.zip","version":"0.11.13"},{"arch":"arm64","filename":"terraform_0.11.13_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_arm64.zip","version":"0.11.13"},{"arch":"386","filename":"terraform_0.11.13_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_openbsd_386.zip","version":"0.11.13"},{"arch":"amd64","filename":"terraform_0.11.13_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_openbsd_amd64.zip","version":"0.11.13"},{"arch":"amd64","filename":"terraform_0.11.13_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_solaris_amd64.zip","version":"0.11.13"},{"arch":"386","filename":"terraform_0.11.13_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_windows_386.zip","version":"0.11.13"},{"arch":"amd64","filename":"terraform_0.11.13_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_windows_amd64.zip","version":"0.11.13"}],"name":"terraform","shasums":"terraform_0.11.13_SHA256SUMS","shasums_signature":"terraform_0.11.13_SHA256SUMS.sig","shasums_signatures":["terraform_0.11.13_SHA256SUMS.72D7468F.sig","terraform_0.11.13_SHA256SUMS.sig"],"version":"0.11.13"},"0.12.0-beta2":{"builds":[{"arch":"amd64","filename":"terraform_0.12.0-beta2_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_darwin_amd64.zip","version":"0.12.0-beta2"},{"arch":"arm64","filename":"terraform_0.12.0-beta2_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_darwin_arm64.zip","version":"0.12.0-beta2"},{"arch":"386","filename":"terraform_0.12.0-beta2_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_freebsd_386.zip","version":"0.12.0-beta2"},{"arch":"amd64","filename":"terraform_0.12.0-beta2_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_freebsd_amd64.zip","version":"0.12.0-beta2"},{"arch":"arm","filename":"terraform_0.12.0-beta2_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_freebsd_arm.zip","version":"0.12.0-beta2"},{"arch":"386","filename":"terraform_0.12.0-beta2_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_linux_386.zip","version":"0.12.0-beta2"},{"arch":"amd64","filename":"terraform_0.12.0-beta2_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_linux_amd64.zip","version":"0.12.0-beta2"},{"arch":"arm","filename":"terraform_0.12.0-beta2_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_linux_arm.zip","version":"0.12.0-beta2"},{"arch":"arm64","filename":"terraform_0.12.0-beta2_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_linux_arm64.zip","version":"0.12.0-beta2"},{"arch":"386","filename":"terraform_0.12.0-beta2_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_openbsd_386.zip","version":"0.12.0-beta2"},{"arch":"amd64","filename":"terraform_0.12.0-beta2_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_openbsd_amd64.zip","version":"0.12.0-beta2"},{"arch":"amd64","filename":"terraform_0.12.0-beta2_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_solaris_amd64.zip","version":"0.12.0-beta2"},{"arch":"386","filename":"terraform_0.12.0-beta2_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_windows_386.zip","version":"0.12.0-beta2"},{"arch":"amd64","filename":"terraform_0.12.0-beta2_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0-beta2/terraform_0.12.0-beta2_windows_amd64.zip","version":"0.12.0-beta2"}],"name":"terraform","shasums":"terraform_0.12.0-beta2_SHA256SUMS","shasums_signature":"terraform_0.12.0-beta2_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.0-beta2_SHA256SUMS.72D7468F.sig","terraform_0.12.0-beta2_SHA256SUMS.sig"],"version":"0.12.0-beta2"},"0.12.0-rc1":{"builds":[{"arch":"amd64","filename":"terraform_0.12.0-rc1_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_darwin_amd64.zip","version":"0.12.0-rc1"},{"arch":"arm64","filename":"terraform_0.12.0-rc1_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_darwin_arm64.zip","version":"0.12.0-rc1"},{"arch":"386","filename":"terraform_0.12.0-rc1_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_freebsd_386.zip","version":"0.12.0-rc1"},{"arch":"amd64","filename":"terraform_0.12.0-rc1_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_freebsd_amd64.zip","version":"0.12.0-rc1"},{"arch":"arm","filename":"terraform_0.12.0-rc1_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_freebsd_arm.zip","version":"0.12.0-rc1"},{"arch":"386","filename":"terraform_0.12.0-rc1_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_linux_386.zip","version":"0.12.0-rc1"},{"arch":"amd64","filename":"terraform_0.12.0-rc1_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_linux_amd64.zip","version":"0.12.0-rc1"},{"arch":"arm","filename":"terraform_0.12.0-rc1_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_linux_arm.zip","version":"0.12.0-rc1"},{"arch":"arm64","filename":"terraform_0.12.0-rc1_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_linux_arm64.zip","version":"0.12.0-rc1"},{"arch":"386","filename":"terraform_0.12.0-rc1_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_openbsd_386.zip","version":"0.12.0-rc1"},{"arch":"amd64","filename":"terraform_0.12.0-rc1_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_openbsd_amd64.zip","version":"0.12.0-rc1"},{"arch":"amd64","filename":"terraform_0.12.0-rc1_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_solaris_amd64.zip","version":"0.12.0-rc1"},{"arch":"386","filename":"terraform_0.12.0-rc1_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_windows_386.zip","version":"0.12.0-rc1"},{"arch":"amd64","filename":"terraform_0.12.0-rc1_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0-rc1/terraform_0.12.0-rc1_windows_amd64.zip","version":"0.12.0-rc1"}],"name":"terraform","shasums":"terraform_0.12.0-rc1_SHA256SUMS","shasums_signature":"terraform_0.12.0-rc1_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.0-rc1_SHA256SUMS.72D7468F.sig","terraform_0.12.0-rc1_SHA256SUMS.sig"],"version":"0.12.0-rc1"},"0.12.0":{"builds":[{"arch":"amd64","filename":"terraform_0.12.0_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_darwin_amd64.zip","version":"0.12.0"},{"arch":"arm64","filename":"terraform_0.12.0_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_darwin_arm64.zip","version":"0.12.0"},{"arch":"386","filename":"terraform_0.12.0_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_freebsd_386.zip","version":"0.12.0"},{"arch":"amd64","filename":"terraform_0.12.0_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_freebsd_amd64.zip","version":"0.12.0"},{"arch":"arm","filename":"terraform_0.12.0_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_freebsd_arm.zip","version":"0.12.0"},{"arch":"386","filename":"terraform_0.12.0_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_linux_386.zip","version":"0.12.0"},{"arch":"amd64","filename":"terraform_0.12.0_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_linux_amd64.zip","version":"0.12.0"},{"arch":"arm","filename":"terraform_0.12.0_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_linux_arm.zip","version":"0.12.0"},{"arch":"arm64","filename":"terraform_0.12.0_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_linux_arm64.zip","version":"0.12.0"},{"arch":"386","filename":"terraform_0.12.0_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_openbsd_386.zip","version":"0.12.0"},{"arch":"amd64","filename":"terraform_0.12.0_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_openbsd_amd64.zip","version":"0.12.0"},{"arch":"amd64","filename":"terraform_0.12.0_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_solaris_amd64.zip","version":"0.12.0"},{"arch":"386","filename":"terraform_0.12.0_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_windows_386.zip","version":"0.12.0"},{"arch":"amd64","filename":"terraform_0.12.0_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_windows_amd64.zip","version":"0.12.0"}],"name":"terraform","shasums":"terraform_0.12.0_SHA256SUMS","shasums_signature":"terraform_0.12.0_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.0_SHA256SUMS.72D7468F.sig","terraform_0.12.0_SHA256SUMS.sig"],"version":"0.12.0"},"0.12.1":{"builds":[{"arch":"amd64","filename":"terraform_0.12.1_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_darwin_amd64.zip","version":"0.12.1"},{"arch":"arm64","filename":"terraform_0.12.1_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_darwin_arm64.zip","version":"0.12.1"},{"arch":"386","filename":"terraform_0.12.1_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_freebsd_386.zip","version":"0.12.1"},{"arch":"amd64","filename":"terraform_0.12.1_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_freebsd_amd64.zip","version":"0.12.1"},{"arch":"arm","filename":"terraform_0.12.1_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_freebsd_arm.zip","version":"0.12.1"},{"arch":"386","filename":"terraform_0.12.1_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_linux_386.zip","version":"0.12.1"},{"arch":"amd64","filename":"terraform_0.12.1_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_linux_amd64.zip","version":"0.12.1"},{"arch":"arm","filename":"terraform_0.12.1_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_linux_arm.zip","version":"0.12.1"},{"arch":"arm64","filename":"terraform_0.12.1_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_linux_arm64.zip","version":"0.12.1"},{"arch":"386","filename":"terraform_0.12.1_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_openbsd_386.zip","version":"0.12.1"},{"arch":"amd64","filename":"terraform_0.12.1_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_openbsd_amd64.zip","version":"0.12.1"},{"arch":"amd64","filename":"terraform_0.12.1_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_solaris_amd64.zip","version":"0.12.1"},{"arch":"386","filename":"terraform_0.12.1_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_windows_386.zip","version":"0.12.1"},{"arch":"amd64","filename":"terraform_0.12.1_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_windows_amd64.zip","version":"0.12.1"}],"name":"terraform","shasums":"terraform_0.12.1_SHA256SUMS","shasums_signature":"terraform_0.12.1_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.1_SHA256SUMS.72D7468F.sig","terraform_0.12.1_SHA256SUMS.sig"],"version":"0.12.1"},"0.12.2":{"builds":[{"arch":"amd64","filename":"terraform_0.12.2_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_darwin_amd64.zip","version":"0.12.2"},{"arch":"arm64","filename":"terraform_0.12.2_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_darwin_arm64.zip","version":"0.12.2"},{"arch":"386","filename":"terraform_0.12.2_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_freebsd_386.zip","version":"0.12.2"},{"arch":"amd64","filename":"terraform_0.12.2_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_freebsd_amd64.zip","version":"0.12.2"},{"arch":"arm","filename":"terraform_0.12.2_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_freebsd_arm.zip","version":"0.12.2"},{"arch":"386","filename":"terraform_0.12.2_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_linux_386.zip","version":"0.12.2"},{"arch":"amd64","filename":"terraform_0.12.2_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_linux_amd64.zip","version":"0.12.2"},{"arch":"arm","filename":"terraform_0.12.2_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_linux_arm.zip","version":"0.12.2"},{"arch":"arm64","filename":"terraform_0.12.2_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_linux_arm64.zip","version":"0.12.2"},{"arch":"386","filename":"terraform_0.12.2_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_openbsd_386.zip","version":"0.12.2"},{"arch":"amd64","filename":"terraform_0.12.2_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_openbsd_amd64.zip","version":"0.12.2"},{"arch":"amd64","filename":"terraform_0.12.2_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_solaris_amd64.zip","version":"0.12.2"},{"arch":"386","filename":"terraform_0.12.2_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_windows_386.zip","version":"0.12.2"},{"arch":"amd64","filename":"terraform_0.12.2_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.2/terraform_0.12.2_windows_amd64.zip","version":"0.12.2"}],"name":"terraform","shasums":"terraform_0.12.2_SHA256SUMS","shasums_signature":"terraform_0.12.2_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.2_SHA256SUMS.72D7468F.sig","terraform_0.12.2_SHA256SUMS.sig"],"version":"0.12.2"},"0.12.3-beta1":{"builds":[{"arch":"amd64","filename":"terraform_0.12.3-beta1_darwin_amd64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_darwin_amd64.zip","version":"0.12.3-beta1"},{"arch":"arm64","filename":"terraform_0.12.3-beta1_darwin_arm64.zip","name":"terraform","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_darwin_arm64.zip","version":"0.12.3-beta1"},{"arch":"386","filename":"terraform_0.12.3-beta1_freebsd_386.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_freebsd_386.zip","version":"0.12.3-beta1"},{"arch":"amd64","filename":"terraform_0.12.3-beta1_freebsd_amd64.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_freebsd_amd64.zip","version":"0.12.3-beta1"},{"arch":"arm","filename":"terraform_0.12.3-beta1_freebsd_arm.zip","name":"terraform","os":"freebsd","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_freebsd_arm.zip","version":"0.12.3-beta1"},{"arch":"386","filename":"terraform_0.12.3-beta1_linux_386.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_linux_386.zip","version":"0.12.3-beta1"},{"arch":"amd64","filename":"terraform_0.12.3-beta1_linux_amd64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_linux_amd64.zip","version":"0.12.3-beta1"},{"arch":"arm","filename":"terraform_0.12.3-beta1_linux_arm.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_linux_arm.zip","version":"0.12.3-beta1"},{"arch":"arm64","filename":"terraform_0.12.3-beta1_linux_arm64.zip","name":"terraform","os":"linux","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_linux_arm64.zip","version":"0.12.3-beta1"},{"arch":"386","filename":"terraform_0.12.3-beta1_openbsd_386.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_openbsd_386.zip","version":"0.12.3-beta1"},{"arch":"amd64","filename":"terraform_0.12.3-beta1_openbsd_amd64.zip","name":"terraform","os":"openbsd","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_openbsd_amd64.zip","version":"0.12.3-beta1"},{"arch":"amd64","filename":"terraform_0.12.3-beta1_solaris_amd64.zip","name":"terraform","os":"solaris","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_solaris_amd64.zip","version":"0.12.3-beta1"},{"arch":"386","filename":"terraform_0.12.3-beta1_windows_386.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_windows_386.zip","version":"0.12.3-beta1"},{"arch":"amd64","filename":"terraform_0.12.3-beta1_windows_amd64.zip","name":"terraform","os":"windows","url":"https://releases.hashicorp.com/terraform/0.12.3-beta1/terraform_0.12.3-beta1_windows_amd64.zip","version":"0.12.3-beta1"}],"name":"terraform","shasums":"terraform_0.12.3-beta1_SHA256SUMS","shasums_signature":"terraform_0.12.3-beta1_SHA256SUMS.sig","shasums_signatures":["terraform_0.12.3-beta1_SHA256SUMS.72D7468F.sig","terraform_0.12.3-beta1_SHA256SUMS.sig"],"version":"0.12.3-beta1"}}}`

	openTofuJSONData = `{"versions":[{"id":"1.6.0-alpha1","files":["tofu_1.6.0-alpha1_386.apk","tofu_1.6.0-alpha1_386.deb","tofu_1.6.0-alpha1_386.rpm","tofu_1.6.0-alpha1_amd64.apk","tofu_1.6.0-alpha1_amd64.deb","tofu_1.6.0-alpha1_amd64.rpm","tofu_1.6.0-alpha1_arm.apk","tofu_1.6.0-alpha1_arm.deb","tofu_1.6.0-alpha1_arm.rpm","tofu_1.6.0-alpha1_arm64.apk","tofu_1.6.0-alpha1_arm64.deb","tofu_1.6.0-alpha1_arm64.rpm","tofu_1.6.0-alpha1_darwin_amd64.zip","tofu_1.6.0-alpha1_darwin_arm64.zip","tofu_1.6.0-alpha1_freebsd_386.zip","tofu_1.6.0-alpha1_freebsd_amd64.zip","tofu_1.6.0-alpha1_freebsd_arm.zip","tofu_1.6.0-alpha1_linux_386.zip","tofu_1.6.0-alpha1_linux_amd64.zip","tofu_1.6.0-alpha1_linux_arm.zip","tofu_1.6.0-alpha1_linux_arm64.zip","tofu_1.6.0-alpha1_openbsd_386.zip","tofu_1.6.0-alpha1_openbsd_amd64.zip","tofu_1.6.0-alpha1_SHA256SUMS","tofu_1.6.0-alpha1_SHA256SUMS.pem","tofu_1.6.0-alpha1_SHA256SUMS.sig","tofu_1.6.0-alpha1_solaris_amd64.zip","tofu_1.6.0-alpha1_windows_386.zip","tofu_1.6.0-alpha1_windows_amd64.zip"]},{"id":"1.6.2","files":["tofu_1.6.2_386.apk","tofu_1.6.2_386.deb","tofu_1.6.2_386.rpm","tofu_1.6.2_amd64.apk","tofu_1.6.2_amd64.deb","tofu_1.6.2_amd64.rpm","tofu_1.6.2_arm.apk","tofu_1.6.2_arm.deb","tofu_1.6.2_arm.rpm","tofu_1.6.2_arm64.apk","tofu_1.6.2_arm64.deb","tofu_1.6.2_arm64.rpm","tofu_1.6.2_darwin_amd64.zip","tofu_1.6.2_darwin_arm64.zip","tofu_1.6.2_freebsd_386.zip","tofu_1.6.2_freebsd_amd64.zip","tofu_1.6.2_freebsd_arm.zip","tofu_1.6.2_linux_386.zip","tofu_1.6.2_linux_amd64.zip","tofu_1.6.2_linux_arm.zip","tofu_1.6.2_linux_arm64.zip","tofu_1.6.2_openbsd_386.zip","tofu_1.6.2_openbsd_amd64.zip","tofu_1.6.2_SHA256SUMS","tofu_1.6.2_SHA256SUMS.pem","tofu_1.6.2_SHA256SUMS.sig","tofu_1.6.2_solaris_amd64.zip","tofu_1.6.2_windows_386.zip","tofu_1.6.2_windows_amd64.zip"]},{"id":"1.7.0-alpha1","files":["tofu_1.7.0-alpha1_386.apk","tofu_1.7.0-alpha1_386.deb","tofu_1.7.0-alpha1_386.rpm","tofu_1.7.0-alpha1_amd64.apk","tofu_1.7.0-alpha1_amd64.deb","tofu_1.7.0-alpha1_amd64.rpm","tofu_1.7.0-alpha1_arm.apk","tofu_1.7.0-alpha1_arm.deb","tofu_1.7.0-alpha1_arm.rpm","tofu_1.7.0-alpha1_arm64.apk","tofu_1.7.0-alpha1_arm64.deb","tofu_1.7.0-alpha1_arm64.rpm","tofu_1.7.0-alpha1_darwin_amd64.zip","tofu_1.7.0-alpha1_darwin_arm64.zip","tofu_1.7.0-alpha1_freebsd_386.zip","tofu_1.7.0-alpha1_freebsd_amd64.zip","tofu_1.7.0-alpha1_freebsd_arm.zip","tofu_1.7.0-alpha1_linux_386.zip","tofu_1.7.0-alpha1_linux_amd64.zip","tofu_1.7.0-alpha1_linux_arm.zip","tofu_1.7.0-alpha1_linux_arm64.zip","tofu_1.7.0-alpha1_openbsd_386.zip","tofu_1.7.0-alpha1_openbsd_amd64.zip","tofu_1.7.0-alpha1_SHA256SUMS","tofu_1.7.0-alpha1_SHA256SUMS.pem","tofu_1.7.0-alpha1_SHA256SUMS.sig","tofu_1.7.0-alpha1_solaris_amd64.zip","tofu_1.7.0-alpha1_windows_386.zip","tofu_1.7.0-alpha1_windows_amd64.zip"]},{"id":"1.7.0-beta1","files":["tofu_1.7.0-beta1_386.apk","tofu_1.7.0-beta1_386.deb","tofu_1.7.0-beta1_386.rpm","tofu_1.7.0-beta1_amd64.apk","tofu_1.7.0-beta1_amd64.deb","tofu_1.7.0-beta1_amd64.rpm","tofu_1.7.0-beta1_arm.apk","tofu_1.7.0-beta1_arm.deb","tofu_1.7.0-beta1_arm.rpm","tofu_1.7.0-beta1_arm64.apk","tofu_1.7.0-beta1_arm64.deb","tofu_1.7.0-beta1_arm64.rpm","tofu_1.7.0-beta1_darwin_amd64.zip","tofu_1.7.0-beta1_darwin_arm64.zip","tofu_1.7.0-beta1_freebsd_386.zip","tofu_1.7.0-beta1_freebsd_amd64.zip","tofu_1.7.0-beta1_freebsd_arm.zip","tofu_1.7.0-beta1_linux_386.zip","tofu_1.7.0-beta1_linux_amd64.zip","tofu_1.7.0-beta1_linux_arm.zip","tofu_1.7.0-beta1_linux_arm64.zip","tofu_1.7.0-beta1_openbsd_386.zip","tofu_1.7.0-beta1_openbsd_amd64.zip","tofu_1.7.0-beta1_SHA256SUMS","tofu_1.7.0-beta1_SHA256SUMS.pem","tofu_1.7.0-beta1_SHA256SUMS.sig","tofu_1.7.0-beta1_solaris_amd64.zip","tofu_1.7.0-beta1_windows_386.zip","tofu_1.7.0-beta1_windows_amd64.zip"]},{"id":"1.7.0-rc1","files":["tofu_1.7.0-rc1_386.apk","tofu_1.7.0-rc1_386.deb","tofu_1.7.0-rc1_386.rpm","tofu_1.7.0-rc1_amd64.apk","tofu_1.7.0-rc1_amd64.deb","tofu_1.7.0-rc1_amd64.rpm","tofu_1.7.0-rc1_arm.apk","tofu_1.7.0-rc1_arm.deb","tofu_1.7.0-rc1_arm.rpm","tofu_1.7.0-rc1_arm64.apk","tofu_1.7.0-rc1_arm64.deb","tofu_1.7.0-rc1_arm64.rpm","tofu_1.7.0-rc1_darwin_amd64.zip","tofu_1.7.0-rc1_darwin_arm64.zip","tofu_1.7.0-rc1_freebsd_386.zip","tofu_1.7.0-rc1_freebsd_amd64.zip","tofu_1.7.0-rc1_freebsd_arm.zip","tofu_1.7.0-rc1_linux_386.zip","tofu_1.7.0-rc1_linux_amd64.zip","tofu_1.7.0-rc1_linux_arm.zip","tofu_1.7.0-rc1_linux_arm64.zip","tofu_1.7.0-rc1_openbsd_386.zip","tofu_1.7.0-rc1_openbsd_amd64.zip","tofu_1.7.0-rc1_SHA256SUMS","tofu_1.7.0-rc1_SHA256SUMS.pem","tofu_1.7.0-rc1_SHA256SUMS.sig","tofu_1.7.0-rc1_solaris_amd64.zip","tofu_1.7.0-rc1_windows_386.zip","tofu_1.7.0-rc1_windows_amd64.zip"]},{"id":"1.7.0","files":["tofu_1.7.0_386.apk","tofu_1.7.0_386.deb","tofu_1.7.0_386.rpm","tofu_1.7.0_amd64.apk","tofu_1.7.0_amd64.deb","tofu_1.7.0_amd64.rpm","tofu_1.7.0_arm.apk","tofu_1.7.0_arm.deb","tofu_1.7.0_arm.rpm","tofu_1.7.0_arm64.apk","tofu_1.7.0_arm64.deb","tofu_1.7.0_arm64.rpm","tofu_1.7.0_darwin_amd64.zip","tofu_1.7.0_darwin_arm64.zip","tofu_1.7.0_freebsd_386.zip","tofu_1.7.0_freebsd_amd64.zip","tofu_1.7.0_freebsd_arm.zip","tofu_1.7.0_linux_386.zip","tofu_1.7.0_linux_amd64.zip","tofu_1.7.0_linux_arm.zip","tofu_1.7.0_linux_arm64.zip","tofu_1.7.0_openbsd_386.zip","tofu_1.7.0_openbsd_amd64.zip","tofu_1.7.0_SHA256SUMS","tofu_1.7.0_SHA256SUMS.pem","tofu_1.7.0_SHA256SUMS.sig","tofu_1.7.0_solaris_amd64.zip","tofu_1.7.0_windows_386.zip","tofu_1.7.0_windows_amd64.zip"]},{"id":"1.7.1-beta1","files":["tofu_1.7.1-beta1_386.apk","tofu_1.7.1-beta1_386.deb","tofu_1.7.1-beta1_386.rpm","tofu_1.7.1-beta1_amd64.apk","tofu_1.7.1-beta1_amd64.deb","tofu_1.7.1-beta1_amd64.rpm","tofu_1.7.1-beta1_arm.apk","tofu_1.7.1-beta1_arm.deb","tofu_1.7.1-beta1_arm.rpm","tofu_1.7.1-beta1_arm64.apk","tofu_1.7.1-beta1_arm64.deb","tofu_1.7.1-beta1_arm64.rpm","tofu_1.7.1-beta1_darwin_amd64.zip","tofu_1.7.1-beta1_darwin_arm64.zip","tofu_1.7.1-beta1_freebsd_386.zip","tofu_1.7.1-beta1_freebsd_amd64.zip","tofu_1.7.1-beta1_freebsd_arm.zip","tofu_1.7.1-beta1_linux_386.zip","tofu_1.7.1-beta1_linux_amd64.zip","tofu_1.7.1-beta1_linux_arm.zip","tofu_1.7.1-beta1_linux_arm64.zip","tofu_1.7.1-beta1_openbsd_386.zip","tofu_1.7.1-beta1_openbsd_amd64.zip","tofu_1.7.1-beta1_SHA256SUMS","tofu_1.7.1-beta1_SHA256SUMS.pem","tofu_1.7.1-beta1_SHA256SUMS.sig","tofu_1.7.1-beta1_solaris_amd64.zip","tofu_1.7.1-beta1_windows_386.zip","tofu_1.7.1-beta1_windows_amd64.zip"]}]}`
)

// TestGetTFList : Get list from hashicorp
func TestGetTFList(t *testing.T) {
	logger = InitLogger("DEBUG")
	product := GetProductById("terraform")
	list, err := getTFList(product, hashiURL, true)
	if err != nil {
		t.Errorf("Error getting list of versions from %q: %v", hashiURL, err)
	}

	val := "0.1.0"
	var exists bool

	if reflect.TypeOf(list).Kind() == reflect.Slice {
		s := reflect.ValueOf(list)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
			}
		}
	}

	if !exists {
		t.Errorf("Not able to find version: %s", val)
	} else {
		t.Log("Write versions exist (expected)")
	}
}

func compareLists(actual []string, expected []string) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("Slices are not equal length: Expected: %v, actual: %v", expected, actual)
	}

	for i, v := range expected {
		if v != actual[i] {
			return fmt.Errorf("Elements are not the same. Expected: %s, actual: %s", v, actual[i])
		}
	}
	return nil
}

type MockListVersionServerConfig struct {
	EnableHashicorpHTML bool
	EnableHashicorpJSON bool
	EnableOpentofuHTML  bool
	EnableOpentofuJSON  bool
}

func getMockListVersionServer(config MockListVersionServerConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/hashicorp/":
			if config.EnableHashicorpHTML {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(hashicorpBody)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.NotFoundHandler().ServeHTTP(w, r)
			}
		case "/terraform/index.json":
			if config.EnableHashicorpJSON {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(hashicorpJSONData)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.NotFoundHandler().ServeHTTP(w, r)
			}
		case "/opentofu/":
			if config.EnableOpentofuHTML {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(openTofuBody)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.NotFoundHandler().ServeHTTP(w, r)
			}
		case "/tofu/api.json":
			if config.EnableOpentofuJSON {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(openTofuJSONData)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.NotFoundHandler().ServeHTTP(w, r)
			}
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}

// TestGetVersionsFromBodyHashicorp :  test hashicorp release body
func TestGetVersionsFromBodyHashicorp(t *testing.T) {
	logger = InitLogger("DEBUG")
	var testTfVersionList tfVersionList
	getVersionsFromBody(hashicorpBody, false, &testTfVersionList)
	expectedVersion := []string{"0.12.2", "0.12.1", "0.12.0", "0.11.13"}
	if err := compareLists(testTfVersionList.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	getVersionsFromBody(hashicorpBody, true, &testTfVersionListPre)
	expectedVersion = []string{"0.12.3-beta1", "0.12.2", "0.12.1", "0.12.0", "0.12.0-rc1", "0.12.0-beta2", "0.11.13"}
	if err := compareLists(testTfVersionListPre.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}
}

// TestGetVersionsFromBodyOpenTofu :  test OpenTofu release body
func TestGetVersionsFromBodyOpenTofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	var testTfVersionList tfVersionList
	getVersionsFromBody(openTofuBody, false, &testTfVersionList)
	expectedVersion := []string{"1.7.0", "1.6.2"}
	if err := compareLists(testTfVersionList.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	getVersionsFromBody(openTofuBody, true, &testTfVersionListPre)
	expectedVersion = []string{"1.7.1-beta1", "1.7.0", "1.7.0-rc1", "1.7.0-beta1", "1.7.0-alpha1", "1.6.2", "1.6.0-alpha1"}
	if err := compareLists(testTfVersionListPre.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}
}

// TestGetVersionFromJSONTerraform
func TestGetVersionFromJSONTerraform(t *testing.T) {
	logger = InitLogger("DEBUG")
	var testTfVersionList tfVersionList
	product := GetProductById("terraform")
	err := getVersionsFromJSON(product, hashicorpJSONData, false, &testTfVersionList)
	assert.NoError(t, err)
	expectedVersion := []string{"0.12.2", "0.12.1", "0.12.0", "0.11.13"}
	if cmpErr := compareLists(testTfVersionList.tflist, expectedVersion); cmpErr != nil {
		t.Errorf("Parsed version does not match expected versions: %v", cmpErr)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	err = getVersionsFromJSON(product, hashicorpJSONData, true, &testTfVersionListPre)
	assert.NoError(t, err)
	expectedVersion = []string{"0.12.3-beta1", "0.12.2", "0.12.1", "0.12.0", "0.12.0-rc1", "0.12.0-beta2", "0.11.13"}
	if cmpErr := compareLists(testTfVersionListPre.tflist, expectedVersion); cmpErr != nil {
		t.Errorf("Parsed pre-release version does not match expected versions: %v", cmpErr)
	}
}

// TestGetVersionFromJSONOpentofu
func TestGetVersionFromJSONOpentofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	var testTfVersionList tfVersionList
	product := GetProductById("opentofu")
	err := getVersionsFromJSON(product, openTofuJSONData, false, &testTfVersionList)
	assert.NoError(t, err)
	expectedVersion := []string{"1.7.0", "1.6.2"}
	if cmpErr := compareLists(testTfVersionList.tflist, expectedVersion); cmpErr != nil {
		t.Errorf("Parsed version does not match expected versions: %v", cmpErr)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	err = getVersionsFromJSON(product, openTofuJSONData, true, &testTfVersionListPre)
	assert.NoError(t, err)
	expectedVersion = []string{"1.7.1-beta1", "1.7.0", "1.7.0-rc1", "1.7.0-beta1", "1.7.0-alpha1", "1.6.2", "1.6.0-alpha1"}
	if cmpErr := compareLists(testTfVersionListPre.tflist, expectedVersion); cmpErr != nil {
		t.Errorf("Parsed pre-release version does not match expected versions: %v", cmpErr)
	}
}

// TestGetTFLatest : Test getTFLatest
func TestGetTFLatest(t *testing.T) {
	logger = InitLogger("DEBUG")
	tests := []struct { // Define a struct for each test case and create a slice of them
		name           string
		product        Product
		serverConfig   MockListVersionServerConfig
		url            string
		expectedLatest string
	}{
		{"Hashicorp JSON", GetProductById("terraform"), MockListVersionServerConfig{EnableHashicorpJSON: true}, "terraform/index.json", "0.12.2"},
		{"Hashicorp List", GetProductById("terraform"), MockListVersionServerConfig{EnableHashicorpHTML: true}, "hashicorp", "0.12.2"},
		{"Opentofu JSON", GetProductById("opentofu"), MockListVersionServerConfig{EnableOpentofuJSON: true}, "tofu/api.json", "1.7.0"},
		{"Opentofu List", GetProductById("opentofu"), MockListVersionServerConfig{EnableOpentofuHTML: true}, "opentofu/", "1.7.0"},
	}

	for test := range slices.Values(tests) {
		t.Run(test.name, func(t *testing.T) {
			server := getMockListVersionServer(test.serverConfig)
			defer server.Close()

			version, err := getTFLatest(test.product, fmt.Sprintf("%s/%s", server.URL, test.url))
			if err != nil {
				t.Error(err)
			}
			if version != test.expectedLatest {
				t.Errorf("Expected latest version does not match. Expected: %s, actual: %s", test.expectedLatest, version)
			}
		})
	}
}

// TestGetTFLatestImplicit : Test getTFLatestImplicit
func TestGetTFLatestImplicit(t *testing.T) {
	logger = InitLogger("DEBUG")
	type versionTest struct {
		version         string
		preRelease      bool
		expectedVersion string
	}
	hashicorpVersions := []versionTest{
		{
			version:         "0.12.0",
			preRelease:      false,
			expectedVersion: "0.12.2",
		},
		{
			version:         "0.11",
			preRelease:      false,
			expectedVersion: "0.12.2",
		},
		{
			version:         "0.12",
			preRelease:      true,
			expectedVersion: "0.12.3-beta1",
		},
	}
	opentofuVersions := []versionTest{
		{
			version:         "1.7.0",
			preRelease:      false,
			expectedVersion: "1.7.0",
		},
		{
			version:         "1.6",
			preRelease:      false,
			expectedVersion: "1.7.0",
		},
		{
			version:         "1.7",
			preRelease:      true,
			expectedVersion: "1.7.1-beta1",
		},
	}
	tests := []struct {
		name         string
		product      Product
		serverConfig MockListVersionServerConfig
		url          string
		versionTests []versionTest
	}{
		{"Hashicorp JSON", GetProductById("terraform"), MockListVersionServerConfig{EnableHashicorpJSON: true}, "terraform/index.json", hashicorpVersions},
		{"Hashicorp List", GetProductById("terraform"), MockListVersionServerConfig{EnableHashicorpHTML: true}, "hashicorp/", hashicorpVersions},
		{"Opentofu JSON", GetProductById("opentofu"), MockListVersionServerConfig{EnableOpentofuJSON: true}, "tofu/api.json", opentofuVersions},
		{"Opentofu List", GetProductById("opentofu"), MockListVersionServerConfig{EnableOpentofuHTML: true}, "opentofu/", opentofuVersions},
	}

	for test := range slices.Values(tests) {
		t.Run(test.name, func(t *testing.T) {
			for _, versionTest := range test.versionTests {
				t.Run(fmt.Sprintf("version=%s,prerelease=%t", versionTest.version, versionTest.preRelease), func(t *testing.T) {
					server := getMockListVersionServer(test.serverConfig)
					defer server.Close()

					version, err := getTFLatestImplicit(test.product, fmt.Sprintf("%s/%s", server.URL, test.url), versionTest.preRelease, versionTest.version)
					if err != nil {
						t.Error(err)
					}
					if version != versionTest.expectedVersion {
						t.Errorf("Expected latest version does not match. Expected: %s, actual: %s", versionTest.expectedVersion, version)
					}
				})
			}
		})
	}
}

// TestGetTFURLBody :  Test getTFURLBody method
func TestGetTFURLBody(t *testing.T) {
	logger = InitLogger("DEBUG")
	server := getMockListVersionServer(MockListVersionServerConfig{EnableHashicorpHTML: true})
	defer server.Close()

	body, err := getTFURLBody(fmt.Sprintf("%s/%s", server.URL, "hashicorp"))
	if err != nil {
		t.Error(err)
	}
	if body != hashicorpBody {
		t.Errorf("Body not returned correctly. Expected: %s, actual: %s", hashicorpBody, body)
	}
}

// TestRemoveDuplicateVersions :  test to removed duplicate
func TestRemoveDuplicateVersions(t *testing.T) {
	logger = InitLogger("DEBUG")
	testArray := []string{"0.0.1", "0.0.2", "0.0.3", "0.0.1", "0.12.0-beta1", "0.12.0-beta1"}

	list := removeDuplicateVersions(testArray)

	if len(list) == len(testArray) {
		t.Errorf("Not able to remove duplicate: %s\n", testArray)
	} else {
		t.Log("Write versions exist (expected)")
	}
}

// TestValidVersionFormat : test if func returns valid version format
func TestValidVersionFormat(t *testing.T) {
	logger = InitLogger("DEBUG")
	var version string
	var valid bool

	version = "0.11.8"
	valid = validVersionFormat(version)
	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9"
	valid = validVersionFormat(version)
	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9-beta1"
	valid = validVersionFormat(version)
	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "0.12.0-rc2"
	valid = validVersionFormat(version)
	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-boom"
	valid = validVersionFormat(version)
	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	// Test valid full version format (using func argument)
	version = "1.11.4"
	valid = validVersionFormat(version, regexSemVer.Full)
	if valid == true {
		t.Logf("Valid full version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify full version format: %s\n", version)
	}

	// Test valid minor version format
	version = "1.11"
	valid = validVersionFormat(version, regexSemVer.Minor)
	if valid == true {
		t.Logf("Valid minor version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify minor version format: %s\n", version)
	}

	// Test valid patch version format
	version = "1.11.4"
	valid = validVersionFormat(version, regexSemVer.Patch)
	if valid == true {
		t.Logf("Valid patch version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify patch version format: %s\n", version)
	}
}

// TestInvalidVersionFormat : test if func catches invalid version format
func TestInvalidVersionFormat(t *testing.T) {
	logger = InitLogger("DEBUG")
	var version string
	var valid bool

	version = "1.11.a"
	valid = validVersionFormat(version)
	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "22323"
	valid = validVersionFormat(version)
	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "@^&*!)!"
	valid = validVersionFormat(version)
	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-01"
	valid = validVersionFormat(version)
	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	// Test invalid full version format (using func argument)
	version = "1.11"
	valid = validVersionFormat(version, regexSemVer.Full)
	if valid == false {
		t.Logf("Invalid full version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify full version format: %s\n", version)
	}

	// Test invalid minor version format
	version = "1.11.4"
	valid = validVersionFormat(version, regexSemVer.Minor)
	if valid == false {
		t.Logf("Invalid minor version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify minor version format: %s\n", version)
	}

	// Test invalid patch version format
	version = "1.11"
	valid = validVersionFormat(version, regexSemVer.Patch)
	if valid == false {
		t.Logf("Invalid patch version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify patch version format: %s\n", version)
	}
}

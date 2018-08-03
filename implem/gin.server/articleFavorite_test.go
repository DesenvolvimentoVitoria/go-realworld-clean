package server_test

import (
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var articleFavoritePath = "/api/articles/" + testData.Article("jane").Slug + "/favorite"

func TestArticleFavoritePost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedComment := testData.Article("")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		FavoritesUpdate(
			testData.User("jane").Name,
			testData.Article("jane").Slug,
			true,
		).
		Return(&expectedComment, nil).
		Times(1)

	jwtHandler := jwt.NewTokenHandler("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	t.Run("happyCase Add to favorites", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(articleFavoritePath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			Expect(t).
			Status(200).
			Done()
	})

	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(articleFavoritePath).
			Expect(t).
			Status(401).
			Done()

	})
}

func TestArticleFavoriteDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedComment := testData.Article("")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		FavoritesUpdate(
			testData.User("jane").Name,
			testData.Article("jane").Slug,
			false,
		).
		Return(&expectedComment, nil).
		Times(1)

	jwtHandler := jwt.NewTokenHandler("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	t.Run("happyCase Remove from to favorites", func(t *testing.T) {
		baloo.New(ts.URL).
			Delete(articleFavoritePath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			Expect(t).
			Status(200).
			JSONSchema(testData.ArticleSingleRespDefinition).
			Done()
	})

	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Delete(articleFavoritePath).
			Expect(t).
			Status(401).
			Done()

	})
}

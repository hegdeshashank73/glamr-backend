package cmd

import (
	"database/sql"
	"io"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/middlewares"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var ginEngine *gin.Engine
var db *sql.DB
var handler *handlers.Handler
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP Server on Port 1729",
	Run: func(cmd *cobra.Command, args []string) {
		// Open Telemetry
		// cleanup := common.InitTracer()
		// defer cleanup(context.Background())
		setupGin()
		setupRoutes()
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func setupGin() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	ginEngine = gin.Default()
	ginEngine.Use(otelgin.Middleware(viper.GetString("OTEL_SERVICE_NAME")))
	ginEngine.Use(middlewares.HandlePanic)
	ginEngine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
	handler = handlers.NewHandler(db)
}

func setupRoutes() {
	authGroup := ginEngine.Group("/")
	authGroup.Use(middlewares.Auth())

	privateGroup := ginEngine.Group("/")
	privateGroup.Use(middlewares.Private(true))

	adminGroup := ginEngine.Group("/")
	adminGroup.Use(middlewares.Admin())

	authNoAuthGroup := ginEngine.Group("/")
	authNoAuthGroup.Use(middlewares.AuthNoAuth())

	authOrPrivateGroup := ginEngine.Group("/")
	authOrPrivateGroup.Use(middlewares.AuthOrPrivate())
	// Health
	ginEngine.GET("/health", handlers.HealthHandler)
	ginEngine.GET("/", handlers.HealthHandler)

	// Metadata
	ginEngine.GET("/test", handlers.TestHandler)
	// ginEngine.GET("/site", handlers.SiteHandler)
	// ginEngine.GET("/site/geo", handlers.SiteGeoHandler)
	// ginEngine.GET("/site/vouches/roles", handlers.VouchRolesHandler)

	// Auth: Magiclinks
	ginEngine.POST("/auth/magiclink", handlers.CreateMagiclinkHandler)
	ginEngine.POST("/auth/magiclink/verify", handlers.VerifyMagiclinkHandler)

	// Auth: Logout
	authGroup.POST("/auth/logout", handlers.LogoutHandler)
	authGroup.POST("/auth/logout/all", handlers.LogoutAllHandler)

	// Availability: Username
	ginEngine.GET("/availability/username", handlers.AvailabilityUsernameHandler)
	adminGroup.POST("/availability/username/:username/reserve", handlers.UsernameReserveHandler)

	// OAuth
	ginEngine.GET("/oauth/:provider/login", handlers.OAuthLoginHandler)
	ginEngine.GET("/oauth/:provider/callback", handlers.OAuthLoginCBHandler)
	ginEngine.POST("/oauth/:provider/callback", handlers.OAuthLoginCBAppHandler)

	//Premium
	authGroup.POST("/intent", handlers.CreateIntentTypeHandler)

	// Webhooks
	ginEngine.POST("/webhooks/payments/razorpay", handlers.WebhookPaymentRazorpayHandler)
	ginEngine.POST("/webhooks/payments/stripe", handlers.WebhookPaymentStripeHandler)

	// Me
	authGroup.GET("/me", handlers.MeHandler)
	authGroup.DELETE("/me", handlers.MeDeleteHandler)

	// Profile: Update
	authGroup.POST("/profile/about", handlers.UpdateProfileAboutHandler)
	authGroup.POST("/profile/socials", handlers.UpdateProfileSocialsHandler)
	authGroup.POST("/profile/domains", handlers.UpdateProfileDomainsHandler)
	authGroup.POST("/profile/img", handlers.UpdateProfileImageHandler)
	authGroup.POST("/profile/booking", handlers.CreatePersonBookingHandler)
	// authGroup.POST("/profile/pin", handlers.PinAnswersHandler)
	// authGroup.GET("/profile/texts/:type", handlers.GetProfileTextsHandler)
	// authGroup.POST("/profile/texts", handlers.CreateTextHandler)
	// authGroup.POST("/profile/texts/:text_id", handlers.UpdateProfileTextHandler)
	// authGroup.DELETE("/profile/texts/:text_id", handlers.DeleteProfileTextHandler)

	// Create Organizations
	authGroup.POST("/organizations", handlers.RegisterOrganizationHandler)
	authNoAuthGroup.GET("/organizations/details", handlers.GetOrganizationByUsernameHandler)
	authNoAuthGroup.GET("/organizations/:organization_id/roles", handlers.GetOrganizationRoleHandler)
	privateGroup.POST("/organizations/:organization_id/img", handlers.UpdateOrganizationImageHandler)
	authOrPrivateGroup.POST("/organizations/:organization_id/roles", handlers.CreateOrganizationRoleHandler)
	authOrPrivateGroup.POST("/organizations/:organization_id/roles/:role_id", handlers.UpdateOrganizationRoleHandler)

	// Applicant APIs
	authGroup.POST("/applications/:application_id/resume", handlers.UpdateApplicationResumeHandler)
	authGroup.GET("/applications/:application_id/resume", handlers.GetApplicationResumeHandler)
	authGroup.POST("/applications/:application_id/answers", handlers.UpdateApplicationAnswersHandler)
	authGroup.POST("/applications/:application_id/freeze", handlers.FreezeApplicationHandler)
	authGroup.GET("/applications", handlers.GetPersonApplicationHandler)
	authGroup.POST("/applications", handlers.CreateAndGetApplicationHandler)
	privateGroup.DELETE("/applications/:application_id", handlers.DeleteApplicationHandler)

	// Recruiter Dashboard APIs
	authGroup.POST("/roles/:role_id/applications/:application_id/actions", handlers.UpdateApplicationActionsHandler)
	authGroup.POST("/roles/:role_id/applications", handlers.GetApplicationsForRoleHandler)

	authGroup.POST("/roles/questions/:question_id/answers", handlers.CreateRoleAnswerHandler)
	authGroup.POST("/roles/answers/:answer_id", handlers.UpdateRoleAnswerHandler)
	authGroup.GET("/roles/questions/:question_id/answers", handlers.GetRoleQuestionAnswersHandler)
	authGroup.POST("/roles/answers/:answer_id/transcription", handlers.UpdateRoleAnswerTranscriptionHandler)
	authGroup.GET("/roles/answers/:answer_id/transcription", handlers.CheckRoleAnswerTranscribedHandler)

	// Experiences
	authGroup.POST("/experiences", handlers.CreateExperienceHandler)
	authGroup.GET("/experiences", handlers.GetExperiencesHandler)
	authGroup.POST("/experiences/:experience_id", handlers.UpdateExperienceHandler)
	authGroup.DELETE("/experiences/:experience_id", handlers.DeleteExperienceHandler)
	privateGroup.POST("/experiences/:experience_id/enrich", handlers.EnrichExperienceHandler)

	// Project
	authGroup.POST("/projects", handlers.CreateProjectHandler)
	authGroup.GET("/projects", handlers.GetProjectsHandler)
	authGroup.POST("/projects/:project_id", handlers.UpdateProjectHandler)
	authGroup.DELETE("/projects/:project_id", handlers.DeleteProjectHandler)

	// Education
	authGroup.POST("/educations", handlers.CreateEducationHandler)
	authGroup.GET("/educations", handlers.GetEducationsHandler)
	authGroup.POST("/educations/:education_id", handlers.UpdateEducationHandler)
	authGroup.DELETE("/educations/:education_id", handlers.DeleteEducationHandler)
	// Patent
	authGroup.POST("/patents", handlers.CreatePatentHandler)
	authGroup.GET("/patents", handlers.GetPatentsHandler)
	authGroup.POST("/patents/:patent_id", handlers.UpdatePatentHandler)
	authGroup.DELETE("/patents/:patent_id", handlers.DeletePatentHandler)
	// Publications
	authGroup.POST("/publications", handlers.CreatePublicationHandler)
	authGroup.GET("/publications", handlers.GetPublicationsHandler)
	authGroup.POST("/publications/:publication_id", handlers.UpdatePublicationHandler)
	authGroup.DELETE("/publications/:publication_id", handlers.DeletePublicationHandler)
	// Certifications
	authGroup.POST("/certifications", handlers.CreateCertificationHandler)
	authGroup.GET("/certifications", handlers.GetCertificationsHandler)
	authGroup.POST("/certifications/:certification_id", handlers.UpdateCertificationHandler)
	authGroup.DELETE("/certifications/:certification_id", handlers.DeleteCertificationHandler)
	//Awards
	authGroup.POST("/awards", handlers.CreateAwardHandler)
	authGroup.GET("/awards", handlers.GetAwardsHandler)
	authGroup.POST("/awards/:award_id", handlers.UpdateAwardHandler)
	authGroup.DELETE("/awards/:award_id", handlers.DeleteAwardHandler)

	// Authorization Requests
	authGroup.GET("/requests", handlers.GetAuthorizationRequestsHandler)
	authGroup.POST("/requests", handlers.CreateAuthorizationRequestHandler)
	authGroup.POST("/requests/:request_id", handlers.UpdateAuthorizationRequestHandler)
	authGroup.DELETE("/requests/:request_id", handlers.DeleteAuthorizationRequestHandler)

	//Referrals, Subscriptions_plans, payment_intent
	authGroup.GET("/refer", handlers.GetReferralCodeHandler)
	ginEngine.GET("/subscriptions/plans", handlers.GetSubscriptionPlansHandler)
	privateGroup.POST("/subscriptions/plans", handlers.CreateSubscriptionPlanHandler)
	authGroup.GET("/subscriptions/plans/:plan_id/checkout", handlers.GetSubscriptionPlanCheckoutHandler)
	authGroup.POST("/subscriptions/plans/:plan_id/checkout", handlers.CreatePaymentIntentHandler)
	authGroup.POST("/subscriptions/verify", handlers.VerifySubscriptionPaymentHandler)
	authGroup.POST("/subscriptions/cancel", handlers.CancelPaymentIntentHandler)

	//subscriptions
	privateGroup.POST("/subscriptions", handlers.CreateSubscriptionHandler)
	authGroup.GET("/subscriptions", handlers.GetSubscriptionsHandler)
	authGroup.DELETE("/subscriptions/:subscription_id", handlers.DeleteSubscriptionHandler)

	authNoAuthGroup.GET("/people/:user_id", handlers.GetPersonHandler)
	authNoAuthGroup.GET("/people/:user_id/answers", handlers.GetPersonAnswersHandler)
	authNoAuthGroup.GET("/people/:user_id/experiences", handlers.GetPersonExperiencesHandler)
	authNoAuthGroup.GET("/people/:user_id/educations", handlers.GetPersonEducationsHandler)
	authNoAuthGroup.GET("/people/:user_id/patents", handlers.GetPersonPatentsHandler)
	authNoAuthGroup.GET("/people/:user_id/publications", handlers.GetPersonPublicationsHandler)
	authNoAuthGroup.GET("/people/:user_id/certifications", handlers.GetPersonCertificationsHandler)
	authNoAuthGroup.GET("/people/:user_id/about", handlers.GetPersonAboutMeHandler)
	authNoAuthGroup.GET("/people/:user_id/timeline", handlers.GetPersonTimelineHandler)
	authNoAuthGroup.GET("/people/:user_id/strengths", handlers.GetPersonStrengthsHandler)
	authNoAuthGroup.POST("/people/timeline/events/:timeline_event_id", handlers.UpdatePersonTimelineEventHandler)
	privateGroup.DELETE("/people/:user_id/reset", handlers.ResetPersonEventsHandler)

	authNoAuthGroup.GET("/people/:user_id/details", handlers.GetPersonByUsernameHandler)
	authNoAuthGroup.GET("/people/:user_id/metas", handlers.GetPersonMetasHandler)
	authNoAuthGroup.GET("/people/:user_id/awards", handlers.GetPersonAwardsHandler)
	authNoAuthGroup.GET("/people/:user_id/projects", handlers.GetPersonProjectsHandler)
	privateGroup.POST("/people/:user_id/metas", handlers.UpdatePersonMetasHandler)
	authNoAuthGroup.GET("/people/:user_id/milestones", handlers.GetPersonMilestonesHandler)

	privateGroup.GET("/private/lookup/people", handlers.LookupPeopleHandler)
	privateGroup.POST("/private/people", handlers.CreatePersonHandler)

	authGroup.GET("/milestones", handlers.GetPersonMilestonesHandler)
	authGroup.POST("/milestones", handlers.CreatePersonMilestoneHandler)
	authGroup.POST("/milestones/:milestone_id", handlers.UpdatePersonMilestoneHandler)
	authGroup.DELETE("/milestones/:milestone_id", handlers.DeletePersonMilestoneHandler)

	authOrPrivateGroup.POST("/people/:user_id/checklists", handlers.UpdatePersonChecklistHandler)

	// authNoAuthGroup.GET("/people/:user_id/highlights", handlers.PeopleHighlightsHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks", handlers.PeopleBookmarksHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/answers", handlers.PeopleBookmarksAnswersHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/links", handlers.PeopleBookmarksLinksHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/topics", handlers.PeopleBookmarksTopicsHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/questions", handlers.PeopleBookmarksQuestionsHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/ticks", handlers.PeopleBookmarksTicksHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/posts", handlers.PeopleBookmarksPostsHandler)
	// authNoAuthGroup.GET("/people/:user_id/bookmarks/articles", handlers.PeopleBookmarksArticlessHandler)
	// authNoAuthGroup.GET("/people/:user_id/share", handlers.PeopleShareHandler)
	// authNoAuthGroup.GET("/people/:user_id/followers", handlers.PeopleFollowersHandler)
	// authNoAuthGroup.GET("/people/:user_id/followings", handlers.PeopleFollowingsHandler)

	// TODO
	authGroup.GET("/people/:user_id/enrollments", handlers.GetPersonEnrollmentsHandler)

	// Answer APIs
	authNoAuthGroup.GET("/answers/:answer_id", handlers.GetAnswerHandler)
	authGroup.POST("/answers/:answer_id", handlers.UpdateAnswerHandler)
	authGroup.DELETE("/answers/:answer_id", handlers.DeleteAnswerHandler)

	// authNoAuthGroup.GET("/answers/:answer_id/share", handlers.ShareAnswerHandler)
	authGroup.POST("/answers/:answer_id/comments", handlers.CreateCommentHandler)
	authNoAuthGroup.GET("/answers/:answer_id/comments", handlers.GetCommentsHandler)
	authGroup.POST("/answers/:answer_id/vouches/invite", handlers.AnswerCreateVouchInviteHandler)
	authGroup.POST("/answers/:answer_id/vouches/:vouch_id", handlers.UpdateAnswerVouchHandler)
	authGroup.POST("/answers/:answer_id/followup", handlers.CreateFollowupQuestionHandler)
	authGroup.POST("/answers/followup/:followup_id", handlers.UpdateFollowupAnswerHandler)

	// Reactions
	authGroup.POST("/answers/:answer_id/reactions/:type", handlers.AnswerReactHandler)
	authGroup.DELETE("/answers/:answer_id/reactions/:type", handlers.AnswerUnreactHandler)

	authGroup.DELETE("/vouches/:vouch_id", handlers.DeleteVouchRequestHandler)
	authGroup.POST("/vouches/:vouch_id", handlers.CreateVouchEntityHandler)
	authGroup.POST("/vouches", handlers.CreateVouchRequestHandler)
	authGroup.GET("/vouches", handlers.GetVouchRequestsHandler)
	authNoAuthGroup.GET("/vouches/preview", handlers.GetVouchPreviewHandler)

	// Strength-finder
	authGroup.GET("/strengths/assessment/invite", handlers.GetStrengthAssessmentInviteHandler)
	authGroup.POST("/strengths/assessment", handlers.CreateSelfStrengthAssessmentHandler)
	authGroup.POST("/strengths/assessment/:strength_id", handlers.UpdateStrengthAssessmentHandler)

	authGroup.POST("/comments/:comment_id", handlers.UpdateCommentHandler)
	authNoAuthGroup.GET("/comments/:comment_id", handlers.GetCommentHandler)
	authGroup.DELETE("/comments/:comment_id", handlers.DeleteCommentHandler)
	authNoAuthGroup.GET("/comments/:comment_id/replies", handlers.GetRepliesHandler)
	authGroup.POST("/comments/:comment_id/replies", handlers.CreateReplyHandler)
	authGroup.POST("/comments/:comment_id/reactions/:type", handlers.ReactCommentHandler)
	authGroup.DELETE("/comments/:comment_id//reactions/:type", handlers.UnreactCommentHandler)

	authGroup.GET("/replies/:reply_id", handlers.GetReplyHandler)
	authGroup.POST("/replies/:reply_id", handlers.UpdateReplyHandler)
	authGroup.DELETE("/replies/:reply_id", handlers.ReplyDeleteHandler)
	authGroup.POST("/replies/:reply_id/reactions/:type", handlers.ReactReplyHandler)
	authGroup.DELETE("/replies/:reply_id/reactions/:type", handlers.UnreactReplyHandler)

	// Questions listing
	authGroup.GET("/questions", handlers.GetQuestionsHandler)
	authGroup.POST("/questions/:question_id/answers", handlers.CreateAnswerHandler)

	// Private APIs
	privateGroup.POST("/questions", handlers.CreateQuestionHandler)

	authNoAuthGroup.GET("/questions/:question_id", handlers.GetQuestionHandler)
	authNoAuthGroup.GET("/questions/:question_id/answers", handlers.GetAnswersByQuestionsHandler)
	authGroup.GET("/questions/:question_id/subquestions", handlers.GetSubquestionsHandler)
	authGroup.POST("/questions/:question_id/subquestions/responses", handlers.CreateSubquestionsResponsesHandler)
	adminGroup.POST("/questions/:question_id/subquestions", handlers.CreateSubquestionsHandler)

	authGroup.GET("/questions/:question_id/feedback", handlers.GetPersonQuestionFeedbacksCountHandler)
	authGroup.POST("/questions/:question_id/feedback", handlers.CreatePersonQuestionFeedbackHandler)
	// Pending
	// authGroup.POST("/questions/:question_id/request", handlers.CreateQuestionRequesteeHandler)
	// authNoAuthGroup.GET("/questions/:question_id/share", handlers.ShareQuestionHandler)
	privateGroup.POST("/questions/:question_id", handlers.UpdateQuestionHandler)

	// Topics
	authGroup.GET("/topics", handlers.GetTopicsHandler)
	authNoAuthGroup.GET("/topics/:topic_id/questions", handlers.GetQuestionsByTopicHandler)
	// Pending optimization
	// authNoAuthGroup.GET("/topics/:topic_id", handlers.TopicDetailsHandler)
	// authNoAuthGroup.GET("/topics/:topic_id/share", handlers.TopicShareHandler)
	// authNoAuthGroup.GET("/topics/:topic_id/details", handlers.TopicDetailsSlugHandler)

	// Professions
	authNoAuthGroup.GET("/professions/fields", handlers.GetFieldsHandler)
	authNoAuthGroup.GET("/professions/fields/:field_id/domains", handlers.GetDomainsHandler)

	authGroup.POST("/affinity/:entity_type/:entity_id", handlers.CreateAffinityHandler)
	authGroup.DELETE("/affinity/:entity_type/:entity_id", handlers.DeleteAffinityHandler)

	// authGroup.GET("/search", handlers.SearchHandler)
	// authGroup.GET("/search/all", handlers.SearchHandler)
	// authGroup.GET("/search/people", handlers.SearchPeopleHandler)
	// authGroup.GET("/search/topics", handlers.SearchTopicsHandler)
	// authGroup.GET("/search/questions", handlers.SearchQuestionsHandler)
	// authGroup.GET("/search/tags", handlers.SearchTagsHandler)
	// authGroup.GET("/search/organizations", handlers.SearchOrganizationsHandler)

	// authGroup.GET("/explore", handlers.ExploreHandler)
	authGroup.GET("/explore/answers", handlers.ExploreAnswersHandler)

	// authGroup.GET("/feed/home/test", handlers.FeedTestHandler)
	// authGroup.GET("/feed/explore", handlers.GetFeedHandler)
	// authGroup.GET("/feed/home", handlers.GetFeedHandler)
	// authGroup.GET("/feed/generate", handlers.GenerateFeedHandler)

	authGroup.GET("/inbox", handlers.GetInboxItemsHandler)
	authGroup.POST("/inbox/:item_id/read", handlers.SetInboxItemReadHandler)

	// authGroup.GET("/inbox/rejected", handlers.InboxGetRejectedHandler)

	// TODO: Deprecated should be removed eventually
	// in favour of /iobox/rejected API
	// authGroup.GET("/inbox/questions/rejected", handlers.InboxGetRejectedHandler)

	// authGroup.GET("/inbox/requests", handlers.InboxGetRequestsHandler)
	// authGroup.GET("/inbox/discourse", handlers.InboxGetDiscourseHandler)
	// authGroup.GET("/inbox/affinity", handlers.InboxGetAffinityHandler)
	// authGroup.GET("/inbox/reactions", handlers.InboxGetReactionsHandler)
	// authGroup.POST("/inbox/items/:item_id/:type", handlers.InboxItemReactHandler)
	// authGroup.POST("/inbox/ack", handlers.InboxAckHandler)

	// authGroup.POST("/suggestions/organizations", handlers.SuggestOrganizationHandler)
	authGroup.POST("/suggestions/questions", handlers.SuggestQuestionHandler)

	// authNoAuthGroup.GET("/bookmarks", handlers.BookmarkStatusHandler)
	// authGroup.POST("/bookmarks/links", handlers.CreateBookmarkLinkHandler)
	// authGroup.DELETE("/bookmarks/links/:link_id", handlers.DeleteBookmarkLinkHandler)
	// authGroup.POST("/bookmarks/links/:link_id/:type", handlers.BookmarkLinkReactHandler)
	// authGroup.DELETE("/bookmarks/links/:link_id/:type", handlers.BookmarkLinkUnreactHandler)

	// authNoAuthGroup.GET("/events", handlers.EventsHandler)
	// authNoAuthGroup.GET("/events/:event_id", handlers.GetEventHandler)
	// authGroup.POST("/events/:event_id/pay", handlers.EventPayHandler)
	// authGroup.GET("/events/:event_id/pricing/breakdown", handlers.EventPricingBreakdownHandler)
	// authNoAuthGroup.POST("/events/:event_id/waitlist", handlers.JoinEventWaitlistHandler)

	// Drafts: Answers
	authGroup.GET("/drafts", handlers.GetDraftsHandler)
	authGroup.POST("/drafts", handlers.CreateDraftHandler)
	authGroup.POST("/drafts/:draft_id", handlers.UpdateDraftHandler)
	authGroup.GET("/drafts/:draft_id", handlers.GetDraftHandler)
	authGroup.DELETE("/drafts/:draft_id", handlers.DeleteDraftHandler)

	// Drafts: Clips
	// authGroup.POST("/drafts/:draft_id/clips", handlers.CreateClipHandler)
	// authGroup.GET("/drafts/:draft_id/clips", handlers.GetClipsHandler)
	// authGroup.DELETE("/clips/:clip_id", handlers.DeleteClipsHandler)

	// authNoAuthGroup.GET("/tags/:tag/:entity", handlers.GetTagEntityHandler)

	// authGroup.POST("/meta/links", handlers.MetaLinksHandler)
	ginEngine.GET("/meta/milestones", handlers.GetMetaMilestonesHandler)
	ginEngine.GET("/meta/vouches/qualities", handlers.GetMetaVouchQualitiesHandler)
	ginEngine.GET("/meta/vouches/relationships", handlers.GetMetaVouchRelationshipsHandler)
	// authGroup.POST("/meta/inbox", handlers.InboxStatusHandler)
	// authNoAuthGroup.GET("/meta/topics", handlers.MetaTopicsHandler)
	// authNoAuthGroup.GET("/meta/tags", handlers.MetaTagsHandler)

	// authNoAuthGroup.GET("/ads", handlers.GetAdsHandler)

	// authNoAuthGroup.GET("/subscriptions/plans", handlers.GetSubscriptionPlansHandler)

	authNoAuthGroup.GET("/offerings", handlers.GetOfferingsHandler)
	authNoAuthGroup.GET("/offerings/:offering_id", handlers.GetOfferingHandler)
	authGroup.GET("/workshops/:workshop_id/resources", handlers.GetWorkshopResourcesHandler)
	authGroup.GET("/enrollments", handlers.GetPersonEnrollmentsHandler)
	authGroup.GET("/enrollments/:workshop_id/invoice", handlers.GetEnrollmentInvoiceHandler)

	authGroup.POST("/assets/upload", handlers.AssetUploadHandler)
	authGroup.POST("/assets/upload/multipart", handlers.AssetMultipartUploadHandler)
	authGroup.POST("/assets/upload/multipart/complete", handlers.CompleteMultipartUploadHandler)
	privateGroup.POST("/templates/email", handlers.CreateEmailTemplateHandler)
	// authGroup.POST("/notifications/fcm", handlers.UpdateNotificationFCMToken)
	// authGroup.GET("/notifications/preferences", handlers.GetNotificationPreferencesHandler)
	// authGroup.POST("/notifications/preferences", handlers.UpdateNotificationPreferencesHandler)
}

func startServer() {
	defer sentry.Recover()
	logrus.Info("Starting the server on :1729")
	ginEngine.Run(":1729")
}

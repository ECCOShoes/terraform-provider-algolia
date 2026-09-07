package provider

import (
	"context"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ECCOShoes/terraform-provider-algolia/internal/client"
)

var (
	_ resource.Resource                = (*indexConfigResource)(nil)
	_ resource.ResourceWithConfigure   = (*indexConfigResource)(nil)
	_ resource.ResourceWithImportState = (*indexConfigResource)(nil)
)

type indexConfigResource struct {
	client *search.APIClient
}

type indexConfigModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	ResetOnDestroy types.Bool `tfsdk:"reset_on_destroy"`

	// Attributes.
	SearchableAttributes    types.List `tfsdk:"searchable_attributes"`
	AttributesForFaceting   types.List `tfsdk:"attributes_for_faceting"`
	UnretrievableAttributes types.List `tfsdk:"unretrievable_attributes"`
	AttributesToRetrieve    types.List `tfsdk:"attributes_to_retrieve"`

	// Ranking.
	Ranking       types.List `tfsdk:"ranking"`
	CustomRanking types.List `tfsdk:"custom_ranking"`

	// Faceting.
	MaxValuesPerFacet types.Int64  `tfsdk:"max_values_per_facet"`
	SortFacetValuesBy types.String `tfsdk:"sort_facet_values_by"`

	// Typo tolerance.
	MinWordSizeFor1Typo       types.Int64  `tfsdk:"min_word_size_for_1_typo"`
	MinWordSizeFor2Typos      types.Int64  `tfsdk:"min_word_size_for_2_typos"`
	TypoTolerance             types.String `tfsdk:"typo_tolerance"`
	AllowTyposOnNumericTokens types.Bool   `tfsdk:"allow_typos_on_numeric_tokens"`

	// Pagination.
	HitsPerPage         types.Int64 `tfsdk:"hits_per_page"`
	PaginationLimitedTo types.Int64 `tfsdk:"pagination_limited_to"`

	// Highlight / snippet.
	AttributesToHighlight             types.List   `tfsdk:"attributes_to_highlight"`
	AttributesToSnippet               types.List   `tfsdk:"attributes_to_snippet"`
	HighlightPreTag                   types.String `tfsdk:"highlight_pre_tag"`
	HighlightPostTag                  types.String `tfsdk:"highlight_post_tag"`
	SnippetEllipsisText               types.String `tfsdk:"snippet_ellipsis_text"`
	RestrictHighlightAndSnippetArrays types.Bool   `tfsdk:"restrict_highlight_and_snippet_arrays"`
}

// NewIndexConfigResource constructs the index_config resource.
func NewIndexConfigResource() resource.Resource {
	return &indexConfigResource{}
}

func (r *indexConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index_config"
}

func (r *indexConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	listOptComputed := func(desc string) schema.ListAttribute {
		return schema.ListAttribute{
			ElementType:   types.StringType,
			Optional:      true,
			Computed:      true,
			Description:   desc,
			PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
		}
	}

	resp.Schema = schema.Schema{
		Description: "Manages the settings (configuration) of an Algolia index. This resource manages " +
			"index settings only; it does not create or delete the index itself. Destroying the " +
			"resource removes it from Terraform state without changing the index unless reset_on_destroy is set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "The index name, used as the resource identifier.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:      true,
				Description:   "Name of the index whose settings are managed. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"reset_on_destroy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "When true, destroying the resource resets the managed settings to Algolia's " +
					"default values instead of only removing the resource from state. Defaults to false.",
			},

			"searchable_attributes":    listOptComputed("Complete list of attributes used for searching."),
			"attributes_for_faceting":  listOptComputed("Attributes used for faceting and filtering."),
			"unretrievable_attributes": listOptComputed("Attributes that can't be retrieved at query time."),
			"attributes_to_retrieve":   listOptComputed("Attributes to include in the API response."),

			"ranking":        listOptComputed("Determines the order in which Algolia returns your results."),
			"custom_ranking": listOptComputed("Attributes to use as custom ranking."),

			"max_values_per_facet": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Maximum number of facet values to return for each facet.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"sort_facet_values_by": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Order in which to retrieve facet values: \"count\" or \"alpha\".",
				Validators:    []validator.String{stringvalidator.OneOf("count", "alpha")},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"min_word_size_for_1_typo": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Minimum number of characters a word must contain to accept matches with one typo.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"min_word_size_for_2_typos": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Minimum number of characters a word must contain to accept matches with two typos.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"typo_tolerance": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether typo tolerance is enabled and how it's applied: \"true\", \"false\", \"min\", or \"strict\".",
				Validators:    []validator.String{stringvalidator.OneOf("true", "false", "min", "strict")},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"allow_typos_on_numeric_tokens": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to allow typos on numbers in the query.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},

			"hits_per_page": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Number of hits per page.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"pagination_limited_to": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "Maximum number of hits accessible through pagination.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},

			"attributes_to_highlight": listOptComputed("Attributes to highlight."),
			"attributes_to_snippet":   listOptComputed("Attributes to snippet, with an optional number of words."),
			"highlight_pre_tag": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "HTML tag inserted before highlighted parts.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"highlight_post_tag": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "HTML tag inserted after highlighted parts.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"snippet_ellipsis_text": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "String used as an ellipsis indicator when a snippet is truncated.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"restrict_highlight_and_snippet_arrays": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to restrict highlighting and snippeting to items that matched the query.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *indexConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (r *indexConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan indexConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applySettings(ctx, plan, &resp.Diagnostics, &resp.State)
}

func (r *indexConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan indexConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applySettings(ctx, plan, &resp.Diagnostics, &resp.State)
}

func (r *indexConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state indexConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.ID.ValueString()
	settings, err := r.client.GetSettings(r.client.NewApiGetSettingsRequest(name))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read index settings", err.Error())
		return
	}

	r.flattenSettings(name, settings, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *indexConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state indexConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Option B (default): just drop from state, leaving the index settings intact.
	if !state.ResetOnDestroy.ValueBool() {
		return
	}

	// Option C: reset the managed settings to Algolia's default values.
	name := state.ID.ValueString()
	updated, err := r.client.SetSettings(r.client.NewApiSetSettingsRequest(name, defaultIndexSettings()))
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Unable to reset index settings", err.Error())
		return
	}
	if _, err := r.client.WaitForTask(name, updated.TaskID); err != nil {
		resp.Diagnostics.AddError("Unable to confirm index settings reset", err.Error())
	}
}

func (r *indexConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// applySettings expands the model, pushes the settings to Algolia, waits for the
// task to publish, then reads the settings back into state.
func (r *indexConfigResource) applySettings(ctx context.Context, plan indexConfigModel, diags *diag.Diagnostics, state *tfsdk.State) {
	name := plan.Name.ValueString()

	settings, d := r.expandSettings(ctx, plan)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	updated, err := r.client.SetSettings(r.client.NewApiSetSettingsRequest(name, settings))
	if err != nil {
		diags.AddError("Unable to set index settings", err.Error())
		return
	}
	if _, err := r.client.WaitForTask(name, updated.TaskID); err != nil {
		diags.AddError("Unable to confirm index settings update", err.Error())
		return
	}

	got, err := r.client.GetSettings(r.client.NewApiGetSettingsRequest(name))
	if err != nil {
		diags.AddError("Unable to read index settings after update", err.Error())
		return
	}

	r.flattenSettings(name, got, &plan)
	diags.Append(state.Set(ctx, &plan)...)
}

func (r *indexConfigResource) expandSettings(ctx context.Context, model indexConfigModel) (*search.IndexSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	searchable, d := expandStringList(ctx, model.SearchableAttributes)
	diags.Append(d...)
	faceting, d := expandStringList(ctx, model.AttributesForFaceting)
	diags.Append(d...)
	unretrievable, d := expandStringList(ctx, model.UnretrievableAttributes)
	diags.Append(d...)
	retrieve, d := expandStringList(ctx, model.AttributesToRetrieve)
	diags.Append(d...)
	ranking, d := expandStringList(ctx, model.Ranking)
	diags.Append(d...)
	customRanking, d := expandStringList(ctx, model.CustomRanking)
	diags.Append(d...)
	toHighlight, d := expandStringList(ctx, model.AttributesToHighlight)
	diags.Append(d...)
	toSnippet, d := expandStringList(ctx, model.AttributesToSnippet)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	settings := &search.IndexSettings{
		SearchableAttributes:              searchable,
		AttributesForFaceting:             faceting,
		UnretrievableAttributes:           unretrievable,
		AttributesToRetrieve:              retrieve,
		Ranking:                           ranking,
		CustomRanking:                     customRanking,
		MaxValuesPerFacet:                 int32Ptr(model.MaxValuesPerFacet),
		SortFacetValuesBy:                 stringPtr(model.SortFacetValuesBy),
		MinWordSizefor1Typo:               int32Ptr(model.MinWordSizeFor1Typo),
		MinWordSizefor2Typos:              int32Ptr(model.MinWordSizeFor2Typos),
		TypoTolerance:                     expandTypoTolerance(model.TypoTolerance),
		AllowTyposOnNumericTokens:         boolPtr(model.AllowTyposOnNumericTokens),
		HitsPerPage:                       int32Ptr(model.HitsPerPage),
		PaginationLimitedTo:               int32Ptr(model.PaginationLimitedTo),
		AttributesToHighlight:             toHighlight,
		AttributesToSnippet:               toSnippet,
		HighlightPreTag:                   stringPtr(model.HighlightPreTag),
		HighlightPostTag:                  stringPtr(model.HighlightPostTag),
		SnippetEllipsisText:               stringPtr(model.SnippetEllipsisText),
		RestrictHighlightAndSnippetArrays: boolPtr(model.RestrictHighlightAndSnippetArrays),
	}
	return settings, diags
}

func (r *indexConfigResource) flattenSettings(name string, s *search.SettingsResponse, model *indexConfigModel) {
	model.ID = types.StringValue(name)
	model.Name = types.StringValue(name)

	model.SearchableAttributes = flattenStringList(s.SearchableAttributes)
	model.AttributesForFaceting = flattenStringList(s.AttributesForFaceting)
	model.UnretrievableAttributes = flattenStringList(s.UnretrievableAttributes)
	model.AttributesToRetrieve = flattenStringList(s.AttributesToRetrieve)

	model.Ranking = flattenStringList(s.Ranking)
	model.CustomRanking = flattenStringList(s.CustomRanking)

	model.MaxValuesPerFacet = int64Value(s.MaxValuesPerFacet)
	model.SortFacetValuesBy = stringValue(s.SortFacetValuesBy)

	model.MinWordSizeFor1Typo = int64Value(s.MinWordSizefor1Typo)
	model.MinWordSizeFor2Typos = int64Value(s.MinWordSizefor2Typos)
	model.TypoTolerance = flattenTypoTolerance(s.TypoTolerance)
	model.AllowTyposOnNumericTokens = boolValue(s.AllowTyposOnNumericTokens)

	model.HitsPerPage = int64Value(s.HitsPerPage)
	model.PaginationLimitedTo = int64Value(s.PaginationLimitedTo)

	model.AttributesToHighlight = flattenStringList(s.AttributesToHighlight)
	model.AttributesToSnippet = flattenStringList(s.AttributesToSnippet)
	model.HighlightPreTag = stringValue(s.HighlightPreTag)
	model.HighlightPostTag = stringValue(s.HighlightPostTag)
	model.SnippetEllipsisText = stringValue(s.SnippetEllipsisText)
	model.RestrictHighlightAndSnippetArrays = boolValue(s.RestrictHighlightAndSnippetArrays)
}

// expandTypoTolerance converts the string form ("true"/"false"/"min"/"strict")
// into the SDK's union type. A null or unknown value yields nil.
func expandTypoTolerance(v types.String) *search.TypoTolerance {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	switch v.ValueString() {
	case "true":
		return search.BoolAsTypoTolerance(true)
	case "false":
		return search.BoolAsTypoTolerance(false)
	case "min":
		return search.TypoToleranceEnumAsTypoTolerance(search.TYPO_TOLERANCE_ENUM_MIN)
	case "strict":
		return search.TypoToleranceEnumAsTypoTolerance(search.TYPO_TOLERANCE_ENUM_STRICT)
	default:
		return nil
	}
}

// flattenTypoTolerance converts the SDK's union type back into the string form.
func flattenTypoTolerance(v *search.TypoTolerance) types.String {
	if v == nil {
		return types.StringNull()
	}
	if v.Bool != nil {
		return types.StringValue(strconv.FormatBool(*v.Bool))
	}
	if v.TypoToleranceEnum != nil {
		return types.StringValue(string(*v.TypoToleranceEnum))
	}
	return types.StringNull()
}

// defaultIndexSettings returns Algolia's documented default values for the
// managed settings. It is used by reset_on_destroy. List settings whose default
// is "all attributes" cannot be expressed through the typed request (the SDK
// omits empty slices), so they are left unchanged.
func defaultIndexSettings() *search.IndexSettings {
	return &search.IndexSettings{
		AttributesToRetrieve:              []string{"*"},
		Ranking:                           []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"},
		MaxValuesPerFacet:                 ptr(int32(100)),
		SortFacetValuesBy:                 ptr("count"),
		MinWordSizefor1Typo:               ptr(int32(4)),
		MinWordSizefor2Typos:              ptr(int32(8)),
		TypoTolerance:                     search.BoolAsTypoTolerance(true),
		AllowTyposOnNumericTokens:         ptr(true),
		HitsPerPage:                       ptr(int32(20)),
		PaginationLimitedTo:               ptr(int32(1000)),
		HighlightPreTag:                   ptr("<em>"),
		HighlightPostTag:                  ptr("</em>"),
		SnippetEllipsisText:               ptr("…"),
		RestrictHighlightAndSnippetArrays: ptr(false),
	}
}

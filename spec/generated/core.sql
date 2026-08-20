-- # Class: BusinessEntity Description: A legal entity. The company scope.
--     * Slot: id
--     * Slot: legal_name
--     * Slot: entity_form Description: LLC, corporation, sole proprietorship, partnership.
--     * Slot: formation_jurisdiction
--     * Slot: formation_date
--     * Slot: company_peak_employment_prior_cy Description: Peak employment across the ENTIRE company in the prior calendar year. Deliberately named with a company scope, because size-based exemptions are company-wide while industry-based exemptions are per-site.
--     * Slot: company_annual_receipts
--     * Slot: affiliation_group_id Description: Affiliate grouping, which size standards aggregate across.
-- # Class: Establishment Description: A single physical site operated by a business entity.
--     * Slot: id
--     * Slot: business_entity
--     * Slot: site_address
--     * Slot: is_unincorporated
--     * Slot: employment_at_establishment
--     * Slot: floor_area_sqft
--     * Slot: opened_on
--     * Slot: closed_on
--     * Slot: BusinessEntity_id Description: Autocreated FK slot
-- # Class: Process Description: An activity unit beneath an establishment. Exists because chemical-process regimes attach codes and quantity thresholds below the site level, and forbid inheriting the site's primary code as a default.
--     * Slot: id
--     * Slot: establishment
--     * Slot: name
--     * Slot: process_type
--     * Slot: is_covered_process
--     * Slot: Establishment_id Description: Autocreated FK slot
-- # Class: Activity Description: A line of business carried on at an establishment.
--     * Slot: id
--     * Slot: establishment
--     * Slot: description
--     * Slot: is_primary Description: Deliberately NOT uniquely constrained. Some programmes speak of "primary industrial activity(ies)" in the plural, and an activity matching a narrative category is primary regardless of its code.
--     * Slot: primacy_basis Description: receipts, headcount, production_rate, operator_designated, narrative_category
--     * Slot: receipts_share
--     * Slot: primacy_justification_text Description: Receipts ordering is stored as advisory guidance with an operator-supplied justification, never as an enforced rule.
--     * Slot: Establishment_id Description: Autocreated FK slot
-- # Class: ProductOffering Description: A product or service offered. A separate demand-based axis; not derivable from the supply-based industry classification.
--     * Slot: id
--     * Slot: establishment
--     * Slot: revenue_share
--     * Slot: is_regulated_substance
--     * Slot: Establishment_id Description: Autocreated FK slot
-- # Class: ChemicalHolding Description: A chemical held at a process. Data class is restricted.
--     * Slot: id
--     * Slot: process
--     * Slot: cas_number
--     * Slot: max_quantity
--     * Slot: uom
--     * Slot: flashpoint_c
--     * Slot: hazard_category
--     * Slot: container_type
--     * Slot: source Description: owner_declared, sds_parsed, or inspection.
--     * Slot: Process_id Description: Autocreated FK slot
-- # Class: EquipmentItem
--     * Slot: id
--     * Slot: process
--     * Slot: equipment_type
--     * Slot: design_capacity
--     * Slot: uom
--     * Slot: Process_id Description: Autocreated FK slot
-- # Class: AttributeDefinition Description: The registry entry for a fact the system may collect or a rule may test. A rule referencing an unregistered attribute, or one at a scope that disagrees with this entry, fails CI. This is what makes scope-confusion bugs unwriteable.
--     * Slot: uri
--     * Slot: label
--     * Slot: scope_unit
--     * Slot: datatype
--     * Slot: unit_dimension
--     * Slot: enum_ref
--     * Slot: collection_method
--     * Slot: data_class
--     * Slot: llm_egress_allowed
-- # Class: Fact Description: An open-world fact carrying the long tail of attributes that do not warrant a typed column.
--     * Slot: id
--     * Slot: subject_ref
--     * Slot: subject_scope
--     * Slot: attribute
--     * Slot: value_typed
--     * Slot: unit
--     * Slot: effective_from
--     * Slot: effective_to
--     * Slot: source
--     * Slot: confidence
-- # Class: Scheme Description: A classification scheme.
--     * Slot: id
--     * Slot: label
--     * Slot: publisher
--     * Slot: is_industry_scheme Description: False for code spaces that look like industry codes but are not, and must never be crosswalked as though they were.
-- # Class: SchemeVersion Description: A dated edition of a scheme. Vintage is part of concept identity.
--     * Slot: id
--     * Slot: scheme
--     * Slot: vintage
--     * Slot: effective_from
--     * Slot: effective_to
-- # Class: Concept Description: A single code within a scheme version. Identity is the triple (scheme, vintage, code) -- never a bare code.
--     * Slot: id Description: Synthetic key. The natural key is (scheme, vintage, code).
--     * Slot: scheme_version
--     * Slot: code
--     * Slot: title
--     * Slot: definition_text
--     * Slot: level
--     * Slot: parent_code
--     * Slot: revision_status
-- # Class: Correspondence Description: A published crosswalk between two scheme versions, as an addressable and versionable object rather than a lookup table.
--     * Slot: id
--     * Slot: source_scheme_version
--     * Slot: target_scheme_version
--     * Slot: publisher
--     * Slot: published_on
--     * Slot: file_provenance
--     * Slot: coverage_notes
-- # Class: ConceptMapping Description: One row of a correspondence.
--     * Slot: id
--     * Slot: correspondence
--     * Slot: source_concept
--     * Slot: target_concept
--     * Slot: match_type
--     * Slot: apportionment_ratio Description: Local extension. XKOS 1.2 defines no such property.
--     * Slot: match_strength Description: Local extension.
--     * Slot: asserted_by
-- # Class: ClassificationAssignment Description: A candidate or confirmed code for a subject. Candidates are presented ranked for human confirmation; only a confirmed assignment may be the code of record.
--     * Slot: id
--     * Slot: subject_ref
--     * Slot: subject_scope
--     * Slot: scheme_version
--     * Slot: code
--     * Slot: rank
--     * Slot: score
--     * Slot: method
--     * Slot: confirmation_state
--     * Slot: is_code_of_record
--     * Slot: provenance
-- # Class: CodeTranslation Description: The result of walking one or more correspondences. The hop path is retained so a translation can be audited rather than trusted.
--     * Slot: id
--     * Slot: assignment
--     * Slot: target_scheme_version
--     * Slot: is_composable Description: True only when every hop is an exactMatch. Close-match chains never auto-compose.
--     * Slot: review_required
-- # Class: RollupRule Description: How a code at a wider scope derives from a narrower one, as data. Encoding this as rows is what lets one regime forbid the roll-up that another requires.
--     * Slot: id
--     * Slot: target_scope
--     * Slot: source_scope
--     * Slot: method
--     * Slot: authority_citation
--     * Slot: regime
--     * Slot: no_default_from_parent
-- # Class: Jurisdiction
--     * Slot: id
--     * Slot: level
--     * Slot: geoid
--     * Slot: label
--     * Slot: parent
--     * Slot: effective_from
--     * Slot: effective_to
-- # Class: JurisdictionProfile Description: A jurisdiction expressed as an inclusion set with exceptions, because "everywhere except one state" is a common real pattern that a flat column cannot represent.
--     * Slot: id
--     * Slot: global_flag
--     * Slot: residency_required
-- # Class: Authority Description: An agency that issues credentials or administers obligations.
--     * Slot: id
--     * Slot: label
--     * Slot: jurisdiction
-- # Class: LegalSource Description: A retrieved, hashed, point-in-time snapshot of authoritative text. Every obligation traces to one of these.
--     * Slot: id
--     * Slot: citation
--     * Slot: source_url
--     * Slot: source_system
--     * Slot: text_hash
--     * Slot: retrieved_at
--     * Slot: as_of_date
--     * Slot: amendment_date
-- # Class: Regime Description: A regulatory programme that generates obligations.
--     * Slot: id
--     * Slot: label
--     * Slot: scope_unit
--     * Slot: authority
-- # Class: Obligation
--     * Slot: id
--     * Slot: regime
--     * Slot: obligation_type
--     * Slot: scope_unit
--     * Slot: trigger_rule_id
--     * Slot: legal_source
--     * Slot: deadline_rule
--     * Slot: recurrence
--     * Slot: non_waivable Description: Survives every exemption in its regime. A non-waivable obligation is surfaced even when the applicant is otherwise exempt.
-- # Class: Determination Description: A reproducible answer. Pinning the rule version, the law date, and a hash of the inputs is what lets the system answer "why did it say that in March".
--     * Slot: id
--     * Slot: subject_ref
--     * Slot: subject_scope
--     * Slot: obligation
--     * Slot: result
--     * Slot: classification_assignment Description: MUST reference an assignment whose confirmation_state is not `unconfirmed`. Enforced by a database CHECK constraint, not by application code.
--     * Slot: rule_version_id
--     * Slot: engine_version
--     * Slot: as_of_law
--     * Slot: input_snapshot_hash
--     * Slot: evidence_tree Description: The predicate tree annotated with truth values and citations.
--     * Slot: determined_at
-- # Class: ProvenanceRecord Description: What produced a surfaced item. No obligation is displayed without one.
--     * Slot: id
--     * Slot: model_id
--     * Slot: prompt_version
--     * Slot: index_version
--     * Slot: rank_at_selection
--     * Slot: confirmed_by
--     * Slot: confirmed_at
--     * Slot: ui_version
-- # Class: Credential Description: A licence, certification, registration, or permit.
--     * Slot: id
--     * Slot: credential_type
--     * Slot: label
--     * Slot: issuing_authority
--     * Slot: jurisdiction_profile
--     * Slot: industry_code_vintage Description: Local extension supplying the vintage CTDL omits.
--     * Slot: legal_source
--     * Slot: estimated_cost
--     * Slot: renewal_period_months
-- # Class: Requirement Description: A reified condition between a credential and what it demands. A table rather than a foreign key, because the edge carries jurisdiction, residency, experience, cost, and dates. Recursive: AND groups contain OR alternative sets contain further groups, which is what real licensing paths need.
--     * Slot: id
--     * Slot: credential
--     * Slot: parent
--     * Slot: node_type
--     * Slot: edge_kind Description: prerequisite, concurrent, conditional-on, or renewal.
--     * Slot: jurisdiction_profile
--     * Slot: residency_required
--     * Slot: min_age
--     * Slot: years_experience
--     * Slot: estimated_cost
--     * Slot: effective_from
--     * Slot: effective_to
--     * Slot: legal_source
--     * Slot: target_credential
--     * Slot: target_predicate Description: A predicate from the rules DSL. This is the seam that unifies the regulatory and credential halves of the system into a single evaluator.
-- # Class: CredentialDependencyEdge Description: A derived, typed edge used for topological ordering. Cycles are a transcription error in the source law and fail CI; they are not a runtime condition to handle.
--     * Slot: id
--     * Slot: from_credential
--     * Slot: to_credential
--     * Slot: edge_kind
--     * Slot: derived_from_requirement
-- # Class: Establishment_jurisdiction_path
--     * Slot: Establishment_id Description: Autocreated FK slot
--     * Slot: jurisdiction_path_id Description: Ordered nation to place. A set, never a scalar.
-- # Class: ClassificationAssignment_alternatives_shown
--     * Slot: ClassificationAssignment_id Description: Autocreated FK slot
--     * Slot: alternatives_shown Description: What the applicant was offered, not merely what they chose. Competent people disagree on the same description; an appeal needs the offered set.
-- # Class: CodeTranslation_result_code
--     * Slot: CodeTranslation_id Description: Autocreated FK slot
--     * Slot: result_code
-- # Class: CodeTranslation_hop_path
--     * Slot: CodeTranslation_id Description: Autocreated FK slot
--     * Slot: hop_path Description: Ordered correspondence identifiers.
-- # Class: CodeTranslation_match_type_chain
--     * Slot: CodeTranslation_id Description: Autocreated FK slot
--     * Slot: match_type_chain
-- # Class: JurisdictionProfile_main_jurisdiction
--     * Slot: JurisdictionProfile_id Description: Autocreated FK slot
--     * Slot: main_jurisdiction_id
-- # Class: JurisdictionProfile_jurisdiction_exception
--     * Slot: JurisdictionProfile_id Description: Autocreated FK slot
--     * Slot: jurisdiction_exception_id
-- # Class: Determination_missing_attributes
--     * Slot: Determination_id Description: Autocreated FK slot
--     * Slot: missing_attributes_uri Description: What must be answered next when the result is indeterminate.
-- # Class: Credential_industry_codes
--     * Slot: Credential_id Description: Autocreated FK slot
--     * Slot: industry_codes Description: Adopted from CTDL, which types this as an untyped string with no vintage validation, so codes from different revisions are indistinguishable.

CREATE TABLE "AttributeDefinition" (
	uri TEXT NOT NULL,
	label TEXT NOT NULL,
	scope_unit VARCHAR(13) NOT NULL,
	datatype TEXT NOT NULL,
	unit_dimension TEXT,
	enum_ref TEXT,
	collection_method TEXT,
	data_class VARCHAR(21) NOT NULL,
	llm_egress_allowed BOOLEAN NOT NULL,
	PRIMARY KEY (uri)
);
CREATE INDEX "ix_AttributeDefinition_uri" ON "AttributeDefinition" (uri);

CREATE TABLE "Scheme" (
	id TEXT NOT NULL,
	label TEXT NOT NULL,
	publisher TEXT,
	is_industry_scheme BOOLEAN,
	PRIMARY KEY (id)
);
CREATE INDEX "ix_Scheme_id" ON "Scheme" (id);

CREATE TABLE "Jurisdiction" (
	id TEXT NOT NULL,
	level VARCHAR(6) NOT NULL,
	geoid TEXT,
	label TEXT NOT NULL,
	parent TEXT,
	effective_from DATE,
	effective_to DATE,
	PRIMARY KEY (id),
	FOREIGN KEY(parent) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_Jurisdiction_id" ON "Jurisdiction" (id);

CREATE TABLE "JurisdictionProfile" (
	id TEXT NOT NULL,
	global_flag BOOLEAN,
	residency_required BOOLEAN,
	PRIMARY KEY (id)
);
CREATE INDEX "ix_JurisdictionProfile_id" ON "JurisdictionProfile" (id);

CREATE TABLE "LegalSource" (
	id TEXT NOT NULL,
	citation TEXT NOT NULL,
	source_url TEXT NOT NULL,
	source_system TEXT,
	text_hash TEXT,
	retrieved_at DATETIME,
	as_of_date DATE,
	amendment_date DATE,
	PRIMARY KEY (id)
);
CREATE INDEX "ix_LegalSource_id" ON "LegalSource" (id);

CREATE TABLE "ProvenanceRecord" (
	id TEXT NOT NULL,
	model_id TEXT,
	prompt_version TEXT,
	index_version TEXT,
	rank_at_selection INTEGER,
	confirmed_by TEXT,
	confirmed_at DATETIME,
	ui_version TEXT,
	PRIMARY KEY (id)
);
CREATE INDEX "ix_ProvenanceRecord_id" ON "ProvenanceRecord" (id);

CREATE TABLE "BusinessEntity" (
	id TEXT NOT NULL,
	legal_name TEXT NOT NULL,
	entity_form TEXT,
	formation_jurisdiction TEXT,
	formation_date DATE,
	company_peak_employment_prior_cy INTEGER,
	company_annual_receipts NUMERIC,
	affiliation_group_id TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(formation_jurisdiction) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_BusinessEntity_id" ON "BusinessEntity" (id);

CREATE TABLE "Fact" (
	id TEXT NOT NULL,
	subject_ref TEXT NOT NULL,
	subject_scope VARCHAR(13) NOT NULL,
	attribute TEXT NOT NULL,
	value_typed TEXT NOT NULL,
	unit TEXT,
	effective_from DATE,
	effective_to DATE,
	source TEXT,
	confidence NUMERIC,
	PRIMARY KEY (id),
	FOREIGN KEY(attribute) REFERENCES "AttributeDefinition" (uri)
);
CREATE INDEX "ix_Fact_id" ON "Fact" (id);

CREATE TABLE "SchemeVersion" (
	id TEXT NOT NULL,
	scheme TEXT NOT NULL,
	vintage TEXT NOT NULL,
	effective_from DATE,
	effective_to DATE,
	PRIMARY KEY (id),
	FOREIGN KEY(scheme) REFERENCES "Scheme" (id)
);
CREATE INDEX "ix_SchemeVersion_id" ON "SchemeVersion" (id);

CREATE TABLE "Authority" (
	id TEXT NOT NULL,
	label TEXT NOT NULL,
	jurisdiction TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(jurisdiction) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_Authority_id" ON "Authority" (id);

CREATE TABLE "JurisdictionProfile_main_jurisdiction" (
	"JurisdictionProfile_id" TEXT,
	main_jurisdiction_id TEXT,
	PRIMARY KEY ("JurisdictionProfile_id", main_jurisdiction_id),
	FOREIGN KEY("JurisdictionProfile_id") REFERENCES "JurisdictionProfile" (id),
	FOREIGN KEY(main_jurisdiction_id) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_JurisdictionProfile_main_jurisdiction_main_jurisdiction_id" ON "JurisdictionProfile_main_jurisdiction" (main_jurisdiction_id);
CREATE INDEX "ix_JurisdictionProfile_main_jurisdiction_JurisdictionProfile_id" ON "JurisdictionProfile_main_jurisdiction" ("JurisdictionProfile_id");

CREATE TABLE "JurisdictionProfile_jurisdiction_exception" (
	"JurisdictionProfile_id" TEXT,
	jurisdiction_exception_id TEXT,
	PRIMARY KEY ("JurisdictionProfile_id", jurisdiction_exception_id),
	FOREIGN KEY("JurisdictionProfile_id") REFERENCES "JurisdictionProfile" (id),
	FOREIGN KEY(jurisdiction_exception_id) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_JurisdictionProfile_jurisdiction_exception_JurisdictionProfile_id" ON "JurisdictionProfile_jurisdiction_exception" ("JurisdictionProfile_id");
CREATE INDEX "ix_JurisdictionProfile_jurisdiction_exception_jurisdiction_exception_id" ON "JurisdictionProfile_jurisdiction_exception" (jurisdiction_exception_id);

CREATE TABLE "Establishment" (
	id TEXT NOT NULL,
	business_entity TEXT NOT NULL,
	site_address TEXT,
	is_unincorporated BOOLEAN,
	employment_at_establishment INTEGER,
	floor_area_sqft INTEGER,
	opened_on DATE,
	closed_on DATE,
	"BusinessEntity_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(business_entity) REFERENCES "BusinessEntity" (id),
	FOREIGN KEY("BusinessEntity_id") REFERENCES "BusinessEntity" (id)
);
CREATE INDEX "ix_Establishment_id" ON "Establishment" (id);

CREATE TABLE "Concept" (
	id TEXT NOT NULL,
	scheme_version TEXT NOT NULL,
	code TEXT NOT NULL,
	title TEXT NOT NULL,
	definition_text TEXT,
	level INTEGER,
	parent_code TEXT,
	revision_status VARCHAR(9),
	PRIMARY KEY (id),
	UNIQUE (scheme_version, code),
	FOREIGN KEY(scheme_version) REFERENCES "SchemeVersion" (id)
);
CREATE INDEX "ix_Concept_id" ON "Concept" (id);
CREATE INDEX "Concept_scheme_version_code_idx" ON "Concept" (scheme_version, code);

CREATE TABLE "Correspondence" (
	id TEXT NOT NULL,
	source_scheme_version TEXT NOT NULL,
	target_scheme_version TEXT NOT NULL,
	publisher TEXT,
	published_on DATE,
	file_provenance TEXT,
	coverage_notes TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(source_scheme_version) REFERENCES "SchemeVersion" (id),
	FOREIGN KEY(target_scheme_version) REFERENCES "SchemeVersion" (id)
);
CREATE INDEX "ix_Correspondence_id" ON "Correspondence" (id);

CREATE TABLE "ClassificationAssignment" (
	id TEXT NOT NULL,
	subject_ref TEXT NOT NULL,
	subject_scope VARCHAR(13) NOT NULL,
	scheme_version TEXT NOT NULL,
	code TEXT NOT NULL,
	rank INTEGER,
	score NUMERIC,
	method VARCHAR(20) NOT NULL,
	confirmation_state VARCHAR(15) NOT NULL,
	is_code_of_record BOOLEAN NOT NULL,
	provenance TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(scheme_version) REFERENCES "SchemeVersion" (id),
	FOREIGN KEY(provenance) REFERENCES "ProvenanceRecord" (id)
);
CREATE INDEX "ix_ClassificationAssignment_id" ON "ClassificationAssignment" (id);

CREATE TABLE "Regime" (
	id TEXT NOT NULL,
	label TEXT NOT NULL,
	scope_unit VARCHAR(13) NOT NULL,
	authority TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(authority) REFERENCES "Authority" (id)
);
CREATE INDEX "ix_Regime_id" ON "Regime" (id);

CREATE TABLE "Credential" (
	id TEXT NOT NULL,
	credential_type VARCHAR(13) NOT NULL,
	label TEXT NOT NULL,
	issuing_authority TEXT NOT NULL,
	jurisdiction_profile TEXT,
	industry_code_vintage TEXT,
	legal_source TEXT,
	estimated_cost NUMERIC,
	renewal_period_months INTEGER,
	PRIMARY KEY (id),
	FOREIGN KEY(issuing_authority) REFERENCES "Authority" (id),
	FOREIGN KEY(jurisdiction_profile) REFERENCES "JurisdictionProfile" (id),
	FOREIGN KEY(legal_source) REFERENCES "LegalSource" (id)
);
CREATE INDEX "ix_Credential_id" ON "Credential" (id);

CREATE TABLE "Process" (
	id TEXT NOT NULL,
	establishment TEXT NOT NULL,
	name TEXT,
	process_type TEXT,
	is_covered_process BOOLEAN,
	"Establishment_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(establishment) REFERENCES "Establishment" (id),
	FOREIGN KEY("Establishment_id") REFERENCES "Establishment" (id)
);
CREATE INDEX "ix_Process_id" ON "Process" (id);

CREATE TABLE "Activity" (
	id TEXT NOT NULL,
	establishment TEXT NOT NULL,
	description TEXT NOT NULL,
	is_primary BOOLEAN,
	primacy_basis TEXT,
	receipts_share NUMERIC,
	primacy_justification_text TEXT,
	"Establishment_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(establishment) REFERENCES "Establishment" (id),
	FOREIGN KEY("Establishment_id") REFERENCES "Establishment" (id)
);
CREATE INDEX "ix_Activity_id" ON "Activity" (id);

CREATE TABLE "ProductOffering" (
	id TEXT NOT NULL,
	establishment TEXT,
	revenue_share NUMERIC,
	is_regulated_substance BOOLEAN,
	"Establishment_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(establishment) REFERENCES "Establishment" (id),
	FOREIGN KEY("Establishment_id") REFERENCES "Establishment" (id)
);
CREATE INDEX "ix_ProductOffering_id" ON "ProductOffering" (id);

CREATE TABLE "ConceptMapping" (
	id TEXT NOT NULL,
	correspondence TEXT NOT NULL,
	source_concept TEXT,
	target_concept TEXT,
	match_type VARCHAR(12) NOT NULL,
	apportionment_ratio NUMERIC,
	match_strength NUMERIC,
	asserted_by TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(correspondence) REFERENCES "Correspondence" (id),
	FOREIGN KEY(source_concept) REFERENCES "Concept" (id),
	FOREIGN KEY(target_concept) REFERENCES "Concept" (id)
);
CREATE INDEX "ix_ConceptMapping_id" ON "ConceptMapping" (id);

CREATE TABLE "CodeTranslation" (
	id TEXT NOT NULL,
	assignment TEXT NOT NULL,
	target_scheme_version TEXT NOT NULL,
	is_composable BOOLEAN,
	review_required BOOLEAN,
	PRIMARY KEY (id),
	FOREIGN KEY(assignment) REFERENCES "ClassificationAssignment" (id),
	FOREIGN KEY(target_scheme_version) REFERENCES "SchemeVersion" (id)
);
CREATE INDEX "ix_CodeTranslation_id" ON "CodeTranslation" (id);

CREATE TABLE "RollupRule" (
	id TEXT NOT NULL,
	target_scope VARCHAR(13) NOT NULL,
	source_scope VARCHAR(13) NOT NULL,
	method VARCHAR(19) NOT NULL,
	authority_citation TEXT NOT NULL,
	regime TEXT,
	no_default_from_parent BOOLEAN,
	PRIMARY KEY (id),
	FOREIGN KEY(regime) REFERENCES "Regime" (id)
);
CREATE INDEX "ix_RollupRule_id" ON "RollupRule" (id);

CREATE TABLE "Obligation" (
	id TEXT NOT NULL,
	regime TEXT NOT NULL,
	obligation_type TEXT NOT NULL,
	scope_unit VARCHAR(13) NOT NULL,
	trigger_rule_id TEXT NOT NULL,
	legal_source TEXT NOT NULL,
	deadline_rule TEXT,
	recurrence TEXT,
	non_waivable BOOLEAN,
	PRIMARY KEY (id),
	FOREIGN KEY(regime) REFERENCES "Regime" (id),
	FOREIGN KEY(legal_source) REFERENCES "LegalSource" (id)
);
CREATE INDEX "ix_Obligation_id" ON "Obligation" (id);

CREATE TABLE "Requirement" (
	id TEXT NOT NULL,
	credential TEXT NOT NULL,
	parent TEXT,
	node_type VARCHAR(9) NOT NULL,
	edge_kind TEXT,
	jurisdiction_profile TEXT,
	residency_required BOOLEAN,
	min_age INTEGER,
	years_experience NUMERIC,
	estimated_cost NUMERIC,
	effective_from DATE,
	effective_to DATE,
	legal_source TEXT NOT NULL,
	target_credential TEXT,
	target_predicate TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(credential) REFERENCES "Credential" (id),
	FOREIGN KEY(parent) REFERENCES "Requirement" (id),
	FOREIGN KEY(jurisdiction_profile) REFERENCES "JurisdictionProfile" (id),
	FOREIGN KEY(legal_source) REFERENCES "LegalSource" (id),
	FOREIGN KEY(target_credential) REFERENCES "Credential" (id)
);
CREATE INDEX "ix_Requirement_id" ON "Requirement" (id);

CREATE TABLE "Establishment_jurisdiction_path" (
	"Establishment_id" TEXT,
	jurisdiction_path_id TEXT,
	PRIMARY KEY ("Establishment_id", jurisdiction_path_id),
	FOREIGN KEY("Establishment_id") REFERENCES "Establishment" (id),
	FOREIGN KEY(jurisdiction_path_id) REFERENCES "Jurisdiction" (id)
);
CREATE INDEX "ix_Establishment_jurisdiction_path_jurisdiction_path_id" ON "Establishment_jurisdiction_path" (jurisdiction_path_id);
CREATE INDEX "ix_Establishment_jurisdiction_path_Establishment_id" ON "Establishment_jurisdiction_path" ("Establishment_id");

CREATE TABLE "ClassificationAssignment_alternatives_shown" (
	"ClassificationAssignment_id" TEXT,
	alternatives_shown TEXT,
	PRIMARY KEY ("ClassificationAssignment_id", alternatives_shown),
	FOREIGN KEY("ClassificationAssignment_id") REFERENCES "ClassificationAssignment" (id)
);
CREATE INDEX "ix_ClassificationAssignment_alternatives_shown_ClassificationAssignment_id" ON "ClassificationAssignment_alternatives_shown" ("ClassificationAssignment_id");
CREATE INDEX "ix_ClassificationAssignment_alternatives_shown_alternatives_shown" ON "ClassificationAssignment_alternatives_shown" (alternatives_shown);

CREATE TABLE "Credential_industry_codes" (
	"Credential_id" TEXT,
	industry_codes TEXT,
	PRIMARY KEY ("Credential_id", industry_codes),
	FOREIGN KEY("Credential_id") REFERENCES "Credential" (id)
);
CREATE INDEX "ix_Credential_industry_codes_industry_codes" ON "Credential_industry_codes" (industry_codes);
CREATE INDEX "ix_Credential_industry_codes_Credential_id" ON "Credential_industry_codes" ("Credential_id");

CREATE TABLE "ChemicalHolding" (
	id TEXT NOT NULL,
	process TEXT NOT NULL,
	cas_number TEXT,
	max_quantity NUMERIC,
	uom TEXT,
	flashpoint_c NUMERIC,
	hazard_category TEXT,
	container_type TEXT,
	source TEXT,
	"Process_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(process) REFERENCES "Process" (id),
	FOREIGN KEY("Process_id") REFERENCES "Process" (id)
);
CREATE INDEX "ix_ChemicalHolding_id" ON "ChemicalHolding" (id);

CREATE TABLE "EquipmentItem" (
	id TEXT NOT NULL,
	process TEXT,
	equipment_type TEXT,
	design_capacity NUMERIC,
	uom TEXT,
	"Process_id" TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(process) REFERENCES "Process" (id),
	FOREIGN KEY("Process_id") REFERENCES "Process" (id)
);
CREATE INDEX "ix_EquipmentItem_id" ON "EquipmentItem" (id);

CREATE TABLE "Determination" (
	id TEXT NOT NULL,
	subject_ref TEXT NOT NULL,
	subject_scope VARCHAR(13) NOT NULL,
	obligation TEXT NOT NULL,
	result VARCHAR(29) NOT NULL,
	classification_assignment TEXT,
	rule_version_id TEXT NOT NULL,
	engine_version TEXT NOT NULL,
	as_of_law DATE NOT NULL,
	input_snapshot_hash TEXT NOT NULL,
	evidence_tree TEXT,
	determined_at DATETIME,
	PRIMARY KEY (id),
	FOREIGN KEY(obligation) REFERENCES "Obligation" (id),
	FOREIGN KEY(classification_assignment) REFERENCES "ClassificationAssignment" (id)
);
CREATE INDEX "ix_Determination_id" ON "Determination" (id);

CREATE TABLE "CredentialDependencyEdge" (
	id TEXT NOT NULL,
	from_credential TEXT NOT NULL,
	to_credential TEXT NOT NULL,
	edge_kind TEXT NOT NULL,
	derived_from_requirement TEXT,
	PRIMARY KEY (id),
	FOREIGN KEY(from_credential) REFERENCES "Credential" (id),
	FOREIGN KEY(to_credential) REFERENCES "Credential" (id),
	FOREIGN KEY(derived_from_requirement) REFERENCES "Requirement" (id)
);
CREATE INDEX "ix_CredentialDependencyEdge_id" ON "CredentialDependencyEdge" (id);

CREATE TABLE "CodeTranslation_result_code" (
	"CodeTranslation_id" TEXT,
	result_code TEXT,
	PRIMARY KEY ("CodeTranslation_id", result_code),
	FOREIGN KEY("CodeTranslation_id") REFERENCES "CodeTranslation" (id)
);
CREATE INDEX "ix_CodeTranslation_result_code_result_code" ON "CodeTranslation_result_code" (result_code);
CREATE INDEX "ix_CodeTranslation_result_code_CodeTranslation_id" ON "CodeTranslation_result_code" ("CodeTranslation_id");

CREATE TABLE "CodeTranslation_hop_path" (
	"CodeTranslation_id" TEXT,
	hop_path TEXT,
	PRIMARY KEY ("CodeTranslation_id", hop_path),
	FOREIGN KEY("CodeTranslation_id") REFERENCES "CodeTranslation" (id)
);
CREATE INDEX "ix_CodeTranslation_hop_path_hop_path" ON "CodeTranslation_hop_path" (hop_path);
CREATE INDEX "ix_CodeTranslation_hop_path_CodeTranslation_id" ON "CodeTranslation_hop_path" ("CodeTranslation_id");

CREATE TABLE "CodeTranslation_match_type_chain" (
	"CodeTranslation_id" TEXT,
	match_type_chain VARCHAR(12),
	PRIMARY KEY ("CodeTranslation_id", match_type_chain),
	FOREIGN KEY("CodeTranslation_id") REFERENCES "CodeTranslation" (id)
);
CREATE INDEX "ix_CodeTranslation_match_type_chain_match_type_chain" ON "CodeTranslation_match_type_chain" (match_type_chain);
CREATE INDEX "ix_CodeTranslation_match_type_chain_CodeTranslation_id" ON "CodeTranslation_match_type_chain" ("CodeTranslation_id");

CREATE TABLE "Determination_missing_attributes" (
	"Determination_id" TEXT,
	missing_attributes_uri TEXT,
	PRIMARY KEY ("Determination_id", missing_attributes_uri),
	FOREIGN KEY("Determination_id") REFERENCES "Determination" (id),
	FOREIGN KEY(missing_attributes_uri) REFERENCES "AttributeDefinition" (uri)
);
CREATE INDEX "ix_Determination_missing_attributes_missing_attributes_uri" ON "Determination_missing_attributes" (missing_attributes_uri);
CREATE INDEX "ix_Determination_missing_attributes_Determination_id" ON "Determination_missing_attributes" ("Determination_id");


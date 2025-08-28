package usecases

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/styles"
	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"go.uber.org/zap"
)

// SimplePDFUseCase defines the interface for simple PDF operations
type SimplePDFUseCase interface {
	GeneratePDF(ctx context.Context, curriculumID uuid.UUID) error
}

// SimplePDFUseCaseImpl is responsible for generating PDFs of curriculums
type SimplePDFUseCaseImpl struct {
	curriculumUseCase CurriculumUseCase
	logger            *zap.Logger
}

// NewSimplePDFUseCase creates a new instance of the PDF generator
func NewSimplePDFUseCase(curriculumUseCase CurriculumUseCase, logger *zap.Logger) SimplePDFUseCase {
	return &SimplePDFUseCaseImpl{
		curriculumUseCase: curriculumUseCase,
		logger:            logger,
	}
}

// GeneratePDF generates the PDF of the curriculum and saves it in the specified paths
func (pg *SimplePDFUseCaseImpl) GeneratePDF(ctx context.Context, curriculumID uuid.UUID) error {
	pg.logger.Info("Starting PDF generation",
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Get the specific curriculum from the database
	curriculum, err := pg.curriculumUseCase.GetCurriculumByID(ctx, curriculumID)
	if err != nil {
		pg.logger.Error("Failed to get curriculum from database",
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("Curriculum retrieved successfully",
		zap.String("curriculum_id", curriculumID.String()),
		zap.String("full_name", curriculum.FullName),
	)

	maroto := pg.createMaroto(*curriculum)

	pg.logger.Debug("Maroto document created, generating PDF",
		zap.String("curriculum_id", curriculumID.String()),
	)

	document, err := maroto.Generate()
	if err != nil {
		pg.logger.Error("Failed to generate PDF document",
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Debug("PDF document generated successfully",
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Create directories if they don't exist
	pdfDir := os.Getenv("PDF_DIR")
	textDir := os.Getenv("TEXT_DIR")

	if err := pg.ensureDirectoryExists(pdfDir); err != nil {
		pg.logger.Error("Failed to create PDF directory",
			zap.String("directory", pdfDir),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}
	if err := pg.ensureDirectoryExists(textDir); err != nil {
		pg.logger.Error("Failed to create text directory",
			zap.String("directory", textDir),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	// Save the PDF
	pdfPath := filepath.Join(pdfDir, "curriculum_simple.pdf")
	err = document.Save(pdfPath)
	if err != nil {
		pg.logger.Error("Failed to save PDF file",
			zap.String("file_path", pdfPath),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("PDF file saved successfully",
		zap.String("file_path", pdfPath),
		zap.String("curriculum_id", curriculumID.String()),
	)

	// Save the report in text
	textPath := filepath.Join(textDir, "curriculum_simple_report.txt")
	err = document.GetReport().Save(textPath)
	if err != nil {
		pg.logger.Error("Failed to save text report",
			zap.String("file_path", textPath),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Error(err),
		)
		return err
	}

	pg.logger.Info("Text report saved successfully",
		zap.String("file_path", textPath),
		zap.String("curriculum_id", curriculumID.String()),
	)

	pg.logger.Info("PDF generation completed successfully",
		zap.String("curriculum_id", curriculumID.String()),
		zap.String("pdf_path", pdfPath),
		zap.String("text_path", textPath),
	)

	return nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (pg *SimplePDFUseCaseImpl) ensureDirectoryExists(dir string) error {
	pg.logger.Debug("Ensuring directory exists",
		zap.String("directory", dir),
	)
	return os.MkdirAll(dir, 0755)
}

// createMaroto creates and configures the Maroto document
func (pg *SimplePDFUseCaseImpl) createMaroto(curriculum dto.CurriculumResponse) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(20).
		WithTopMargin(20).
		WithRightMargin(20).
		WithBottomMargin(15).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	// Add all sections of the curriculum
	m.AddRows(pg.getPersonalDataSection(curriculum)...)
	m.AddRows(pg.getIntroSection(curriculum)...)
	m.AddRows(pg.getExperienceSection(curriculum.Works)...)
	m.AddRows(pg.getEducationSection(curriculum)...)
	m.AddRows(pg.getCoursesSection(curriculum)...)
	m.AddRows(pg.getTechnologiesSection(curriculum)...)
	m.AddRows(pg.getLanguagesSection(curriculum)...)
	m.AddRows(pg.getDriveLicenseSection(curriculum)...)
	m.AddRows(pg.getSocialLinksSection(curriculum)...)

	return m
}

// getPersonalDataSection returns the personal data section
func (pg *SimplePDFUseCaseImpl) getPersonalDataSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		// Full name
		row.New(15).Add(
			text.NewCol(12, curriculum.FullName, props.Text{
				Top:   6,
				Size:  28,
				Style: fontstyle.Bold,
				Align: align.Center,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		// Professional job
		row.New(10).Add(
			text.NewCol(12, curriculum.JobDescription, props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Normal,
				Align: align.Center,
				Color: styles.GetSecondaryColorSimplePDF(),
			}),
		),

		// Contact information
		row.New(8).Add(
			text.NewCol(12, curriculum.Email+" • "+curriculum.Phone, props.Text{
				Top:   4,
				Size:  11,
				Align: align.Center,
				Color: styles.GetTextColorSimplePDF(),
			}),
		),

		row.New(10),
	}
}

// getIntroSection returns the introduction section
func (pg *SimplePDFUseCaseImpl) getIntroSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "About", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(15).Add(
			text.NewCol(12, curriculum.Intro, props.Text{
				Size:  11,
				Align: align.Justify,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}

// getExperienceSection returns the professional experience section
func (pg *SimplePDFUseCaseImpl) getExperienceSection(works []dto.WorkResponse) []core.Row {
	var rows []core.Row

	rows = append(rows, row.New(10).Add(
		text.NewCol(12, "Professional Experience", props.Text{
			Top:   2,
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: styles.GetPrimaryColorSimplePDF(),
		}),
	))

	for i, work := range works {
		// Job title and company
		endDateStr := "Present"
		if work.EndDate != nil {
			endDateStr = work.EndDate.Format("2006-01")
		}

		rows = append(rows, row.New(8).Add(
			text.NewCol(8, work.JobTitle, props.Text{
				Size:  13,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetSecondaryColorSimplePDF(),
				Top:   2,
			}),
			text.NewCol(4, work.StartDate.Format("2006-01")+" - "+endDateStr, props.Text{
				Size:  10,
				Style: fontstyle.Italic,
				Align: align.Right,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		))

		// Company name
		rows = append(rows, row.New(4).Add(
			text.NewCol(12, work.CompanyName, props.Text{
				Size:  12,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
				Top:   1,
			}),
		))

		// Company description
		rows = append(rows, row.New(12).Add(
			text.NewCol(12, work.CompanyDescription, props.Text{
				Size:  10,
				Align: align.Justify,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		))

		if i < len(works)-1 {
			rows = append(rows, row.New(8))
		}
	}

	rows = append(rows, row.New(10))
	return rows
}

// getEducationSection returns the education section
func (pg *SimplePDFUseCaseImpl) getEducationSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Education", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.LevelEducation, props.Text{
				Size:  12,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}

// getCoursesSection returns the courses section
func (pg *SimplePDFUseCaseImpl) getCoursesSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Courses", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.Courses, props.Text{
				Size:  11,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(8),
	}
}

// getTechnologiesSection returns the technologies section
func (pg *SimplePDFUseCaseImpl) getTechnologiesSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Technologies", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.Technologies, props.Text{
				Size:  11,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}

// getLanguagesSection returns the languages section
func (pg *SimplePDFUseCaseImpl) getLanguagesSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Languages", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.Languages, props.Text{
				Size:  11,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}

// getDriveLicenseSection returns the driver's license section
func (pg *SimplePDFUseCaseImpl) getDriveLicenseSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Driver's License", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.DriverLicense, props.Text{
				Size:  11,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}

// getSocialLinksSection returns the social links section
func (pg *SimplePDFUseCaseImpl) getSocialLinksSection(curriculum dto.CurriculumResponse) []core.Row {
	return []core.Row{
		row.New(10).Add(
			text.NewCol(12, "Social Links", props.Text{
				Top:   2,
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Left,
				Color: styles.GetPrimaryColorSimplePDF(),
			}),
		),

		row.New(8).Add(
			text.NewCol(12, curriculum.SocialLinks, props.Text{
				Size:  11,
				Align: align.Left,
				Color: styles.GetTextColorSimplePDF(),
				Top:   2,
			}),
		),

		row.New(10),
	}
}
